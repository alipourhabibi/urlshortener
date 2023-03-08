package ports

import (
	aggregate "github.com/alipourhabibi/urlshortener/internal/core/aggergate"
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type URLService interface {
	Add(string, *entity.User) (string, error)
	Get(string) (aggregate.URL, error)
}

type AuthenticationService interface {
	Register(string, string) (*entity.TokenDetails, error)
	LogIn(string, string) (*entity.TokenDetails, error)
}

type AuthorizationService interface {
	CreateTokensAndMetaData(uuid.UUID) (*entity.TokenDetails, error)
	VerifyToken(string, string) (*jwt.Token, error)
}
