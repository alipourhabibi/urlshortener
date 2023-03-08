package authorization

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alipourhabibi/urlshortener/config"
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/alipourhabibi/urlshortener/internal/core/messages"
	"github.com/alipourhabibi/urlshortener/internal/core/ports"
	redisdb "github.com/alipourhabibi/urlshortener/internal/repository/redis"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JWTConfiguration func(*JWTService) error

type JWTService struct {
	cacheRepoistory ports.CacheRepository
}

type Claims struct {
	TD entity.TokenDetails
	jwt.StandardClaims
}

func New(cfgs ...JWTConfiguration) (*JWTService, error) {
	jc := &JWTService{}

	// Loop through all the cfgs and apply them
	for _, cfg := range cfgs {
		err := cfg(jc)
		if err != nil {
			return nil, err
		}
	}
	return jc, nil
}

func WithRedisCacheRepository(config config.Redis) JWTConfiguration {
	return func(js *JWTService) error {
		client, err := redisdb.New(config)
		if err != nil {
			return err
		}
		js.cacheRepoistory = client
		return nil
	}
}

// Generates Only one JWT based on the given secret and expiration time
func (j *JWTService) generateToken(username string) (*entity.TokenDetails, error) {

	var err error
	td := &entity.TokenDetails{}

	accessExp := time.Now().Add(time.Minute * 30).Unix()
	td.ATExpires = accessExp
	td.AccessTokenUuid = uuid.New().String()

	refreshExp := time.Now().Add(time.Hour * 30 * 7).Unix()
	td.RTExpires = refreshExp
	td.RefreshTokenUuid = td.AccessTokenUuid + "++" + username

	// TODO get configs from dependency injection
	secretKeyAcess := []byte(config.Confs.Auth.AccessToken)
	secretKeyRefresh := []byte(config.Confs.Auth.RefreshToken)

	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.AccessTokenUuid
	atClaims["username"] = username
	atClaims["exp"] = td.ATExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString(secretKeyAcess)
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshTokenUuid
	rtClaims["username"] = username
	rtClaims["exp"] = td.RTExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString(secretKeyRefresh)
	if err != nil {
		return nil, err
	}

	return td, nil

}

// Verify JWT based on given secret which can be for accesstoken and refreshtoken
func (j *JWTService) VerifyToken(tokenString, secret string) (*jwt.Token, error) {
	claims := jwt.MapClaims{}
	secretKey := []byte(secret)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Internal Server Erorr")
	}
	exp, _ := claims["exp"].(float64)
	if token.Valid && time.Unix(int64(exp), 0).Unix() > time.Now().Unix() {
		return token, nil
	}
	return nil, err
}

// Generate Both access token and refresh token and set them in redis as well
func (j *JWTService) CreateTokensAndMetaData(username string) (*entity.TokenDetails, error) {

	td, err := j.generateToken(username)

	if err != nil {
		return nil, err
	}

	// save tokens metadata to redis
	at := time.Unix(td.ATExpires, 0)
	rt := time.Unix(td.RTExpires, 0)
	now := time.Now()

	ATCreated, err := j.cacheRepoistory.Set(td.AccessTokenUuid, username, at.Sub(now))
	if err != nil {
		return nil, err
	}
	RTCreated, err := j.cacheRepoistory.Set(td.RefreshTokenUuid, username, rt.Sub(now))
	if err != nil {
		return nil, err
	}
	if ATCreated == "0" || RTCreated == "0" {
		return nil, fmt.Errorf("no record inserted")
	}

	return td, nil
}

func (r *JWTService) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": messages.ErrUnauthorized.Error(),
			})
			c.Abort()
			return
		}
		parts := strings.Split(authorization, "Bearer")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": messages.ErrBadRequest.Error(),
			})
			c.Abort()
			return
		}
		token := strings.TrimSpace(parts[1])
		if len(token) < 1 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": messages.ErrBadRequest.Error(),
			})
			c.Abort()
			return
		}
		// TODO get config from dependency injection
		t, err := r.VerifyToken(token, config.Confs.Auth.AccessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": messages.ErrUnauthorized.Error(),
			})
			c.Abort()
			return
		}
		claims := t.Claims.(jwt.MapClaims)
		username, ok := claims["username"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg": messages.ErrInternalServerError.Error(),
			})
			c.Abort()
			return
		}
		if username == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": messages.ErrUnauthorized.Error(),
			})
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Next()
	}
}
