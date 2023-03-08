package urlhdl

import (
	"github.com/alipourhabibi/urlshortener/config"
	"github.com/alipourhabibi/urlshortener/internal/core/services/authentication"
	"github.com/alipourhabibi/urlshortener/internal/core/services/url"
	"github.com/gin-gonic/gin"
)

func Launch() error {
	urlService, err := url.New(url.WithPostgresURLRepository(config.Confs.PostgresDB), url.WithRedisCacheRepository(config.Confs.Redis))
	if err != nil {
		return err
	}
	authenticationService, err := authentication.New(authentication.WithJWTService(config.Confs.RedisJWT), authentication.WithAuthRepository(config.Confs.PostgresDB))
	if err != nil {
		return err
	}

	urlHandler := NewUrlHttpHandler(urlService)
	authHandler := NewAuthetnticationHandler(authenticationService)

	r := gin.Default()
	r.GET("api/v1/url/:url", urlHandler.Get)
	r.POST("api/v1/url", urlHandler.Add)

	r.POST("api/v1/register/", authHandler.Register)
	r.POST("api/v1/login", authHandler.LogIn)

	return r.Run(":" + config.Confs.Port)
}
