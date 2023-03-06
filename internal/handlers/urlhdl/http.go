package urlhdl

import (
	"github.com/alipourhabibi/urlshortener/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type UrlHttpHandler struct {
	urlHandler ports.URLService
}

func NewUrlHttpHandler(urlService ports.URLService) *UrlHttpHandler {
	return &UrlHttpHandler{
		urlHandler: urlService,
	}
}

func (hdl *UrlHttpHandler) Get(c *gin.Context) {
}
func (hdl *UrlHttpHandler) Add(c *gin.Context) {
}
