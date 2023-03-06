package urlhdl

import (
	"github.com/alipourhabibi/urlshortener/config"
	"github.com/alipourhabibi/urlshortener/internal/core/services/url"
	"github.com/gin-gonic/gin"
)

func Launch() error {
	urlService, err := url.New(url.WithPostgresURLRepository(config.Confs.PostgresDB))
	if err != nil {
		return err
	}
	urlHandler := NewUrlHttpHandler(urlService)

	r := gin.Default()
	r.GET("api/v1/url/:url", urlHandler.Get)
	r.POST("api/v1/url", urlHandler.Add)

	return r.Run(":" + config.Confs.Port)
}
