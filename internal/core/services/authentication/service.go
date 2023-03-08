package authentication

import (
	"github.com/alipourhabibi/urlshortener/config"
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/alipourhabibi/urlshortener/internal/core/messages"
	"github.com/alipourhabibi/urlshortener/internal/core/ports"
	"github.com/alipourhabibi/urlshortener/internal/core/services/authorization"
	"github.com/alipourhabibi/urlshortener/internal/repository/memory"
	"github.com/alipourhabibi/urlshortener/internal/repository/postgres"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationConfiguration func(*AuthenticationService) error

type AuthenticationService struct {
	authorizationService ports.AuthorizationService
	authRepository       ports.AuthenticationRepository
}

func New(cfgs ...AuthenticationConfiguration) (*AuthenticationService, error) {
	ac := &AuthenticationService{}

	// Loop through all the cfgs and apply them
	for _, cfg := range cfgs {
		err := cfg(ac)
		if err != nil {
			return nil, err
		}
	}
	return ac, nil
}

func WithAuthMemoryRepository() AuthenticationConfiguration {
	return func(as *AuthenticationService) error {
		client := memory.NewMemUser()
		as.authRepository = client
		return nil
	}
}

func WithAuthRepository(config config.Posgres) AuthenticationConfiguration {
	return func(as *AuthenticationService) error {
		client, err := postgres.NewAuthRepository(config)
		if err != nil {
			return err
		}
		as.authRepository = client
		return nil
	}
}

func WithJWTService(redis config.Redis) AuthenticationConfiguration {
	return func(as *AuthenticationService) error {
		service, err := authorization.New(authorization.WithRedisCacheRepository(redis))
		if err != nil {
			return err
		}
		as.authorizationService = service
		return nil
	}
}

func (as *AuthenticationService) Register(username, pass string) (*entity.TokenDetails, error) {
	if username == "" || pass == "" {
		return nil, messages.ErrBadRequest
	}
	exists, err := as.authRepository.Exists(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, messages.ErrAlreadyExists
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		return nil, err
	}
	user := entity.User{
		ID:       uuid.New(),
		Username: username,
		Password: string(hashed),
	}
	err = as.authRepository.Add(user)
	if err != nil {
		return nil, err
	}
	return as.authorizationService.CreateTokensAndMetaData(user.Username)
}

func (as *AuthenticationService) LogIn(username, pass string) (*entity.TokenDetails, error) {
	user, err := as.authRepository.Get(username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if err != nil {
		return nil, messages.ErrUnauthorized
	}
	return as.authorizationService.CreateTokensAndMetaData(user.Username)
}
