package urlhdl

import (
	"errors"
	"net/http"

	"github.com/alipourhabibi/urlshortener/internal/core/messages"
	"github.com/alipourhabibi/urlshortener/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type AuthenticationHandler struct {
	authService ports.AuthenticationService
}

func NewAuthetnticationHandler(authService ports.AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{
		authService: authService,
	}
}

func (ahdl *AuthenticationHandler) Register(c *gin.Context) {
	payload := struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": ErrBadRequest.Error(),
		})
		return
	}
	td, err := ahdl.authService.Register(payload.Username, payload.Password)
	if err != nil {
		if errors.Is(err, messages.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":           "OK",
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	})
}
func (ahdl *AuthenticationHandler) LogIn(c *gin.Context) {
	payload := struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": ErrBadRequest.Error(),
		})
		return
	}
	td, err := ahdl.authService.LogIn(payload.Username, payload.Password)
	if err != nil {
		if errors.Is(err, messages.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":           "OK",
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	})
}
