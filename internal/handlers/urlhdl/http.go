package urlhdl

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alipourhabibi/urlshortener/internal/core/ports"
	"github.com/gin-gonic/gin"
)

var (
	ErrBadRequest          = fmt.Errorf("Bad Request")
	ErrInternalServerError = fmt.Errorf("Internal Server Error")
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
	url := c.Param("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": ErrBadRequest.Error(),
		})
		return
	}
	u, err := hdl.urlHandler.Get(url)
	if err != nil {
		log.Printf("\nInternal Server Error %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": ErrInternalServerError.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":      "OK",
		"original": u.GetOriginal(),
	})
}
func (hdl *UrlHttpHandler) Add(c *gin.Context) {
	payload := struct {
		Url string `json:"url" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": ErrBadRequest.Error(),
		})
		return
	}
	shortened, err := hdl.urlHandler.Add(payload.Url, nil)
	if err != nil {
		log.Printf("\nInternal Server Error %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": ErrInternalServerError.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":       "OK",
		"shortened": shortened,
	})
}
