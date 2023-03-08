package urlhdl

import (
	"github.com/alipourhabibi/urlshortener/config"
	"github.com/alipourhabibi/urlshortener/internal/core/services/authentication"
	"github.com/alipourhabibi/urlshortener/internal/core/services/authorization"
	"github.com/alipourhabibi/urlshortener/internal/core/services/url"
	"github.com/gin-gonic/gin"
)

func Launch() error {
	urlService, err := url.New(url.WithPostgresURLRepository(config.Confs.PostgresDB), url.WithRedisCacheRepository(config.Confs.Redis))
	if err != nil {
		return err
	}
	jwtService := authentication.WithJWTService(config.Confs.RedisJWT)
	authenticationService, err := authentication.New(jwtService, authentication.WithAuthRepository(config.Confs.PostgresDB))
	if err != nil {
		return err
	}
	authorizationService, err := authorization.New(authorization.WithRedisCacheRepository(config.Confs.RedisJWT))
	if err != nil {
		return err
	}

	urlHandler := NewUrlHttpHandler(urlService)
	authHandler := NewAuthetnticationHandler(authenticationService)

	r := gin.Default()
	r.GET("api/v1/url/:url", urlHandler.Get)
	r.GET(":url", urlHandler.Redirect)
	r.POST("api/v1/url", authorizationService.Middleware(), urlHandler.Add)

	r.POST("api/v1/register/", authHandler.Register)
	r.POST("api/v1/login", authHandler.LogIn)

	return r.Run(":" + config.Confs.Port)
}
