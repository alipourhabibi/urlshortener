package authorization

import (
	"fmt"
	"time"

	"github.com/alipourhabibi/urlshortener/config"
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/alipourhabibi/urlshortener/internal/core/ports"
	redisdb "github.com/alipourhabibi/urlshortener/internal/repository/redis"
	"github.com/dgrijalva/jwt-go"
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
func (j *JWTService) generateToken(id uuid.UUID) (*entity.TokenDetails, error) {

	var err error
	td := &entity.TokenDetails{}

	accessExp := time.Now().Add(time.Minute * 30).Unix()
	td.ATExpires = accessExp
	td.AccessTokenUuid = uuid.New().String()

	refreshExp := time.Now().Add(time.Hour * 30 * 7).Unix()
	td.RTExpires = refreshExp
	td.RefreshTokenUuid = td.AccessTokenUuid + "++" + id.String()

	secretKeyAcess := []byte(config.Confs.Auth.AccessToken)
	secretKeyRefresh := []byte(config.Confs.Auth.RefreshToken)

	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.AccessTokenUuid
	atClaims["uuid"] = id
	atClaims["exp"] = td.ATExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString(secretKeyAcess)
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshTokenUuid
	rtClaims["uuid"] = id
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
func (j *JWTService) CreateTokensAndMetaData(id uuid.UUID) (*entity.TokenDetails, error) {

	td, err := j.generateToken(id)

	if err != nil {
		return nil, err
	}

	// save tokens metadata to redis
	at := time.Unix(td.ATExpires, 0)
	rt := time.Unix(td.RTExpires, 0)
	now := time.Now()

	ATCreated, err := j.cacheRepoistory.Set(td.AccessTokenUuid, id, at.Sub(now))
	if err != nil {
		return nil, err
	}
	RTCreated, err := j.cacheRepoistory.Set(td.RefreshTokenUuid, id, rt.Sub(now))
	if err != nil {
		return nil, err
	}
	if ATCreated == "0" || RTCreated == "0" {
		return nil, fmt.Errorf("no record inserted")
	}

	return td, nil
}
