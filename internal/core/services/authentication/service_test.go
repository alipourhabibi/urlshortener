package authentication_test

import (
	"errors"
	"testing"

	"github.com/alipourhabibi/urlshortener/config"
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/alipourhabibi/urlshortener/internal/core/messages"
	"github.com/alipourhabibi/urlshortener/internal/core/services/authentication"
)

func TestAuthentication(t *testing.T) {
	err := config.Confs.Load("./local_config.yaml")
	if err != nil {
		t.Error(err)
	}
	jwtService := authentication.WithJWTService(config.Confs.RedisJWT)
	authenticationService, err := authentication.New(jwtService, authentication.WithAuthMemoryRepository())
	if err != nil {
		t.Error(err)
	}

	user1 := entity.User{
		Username: "ali",
		Password: "ali",
	}
	user2 := entity.User{
		Username: "ali",
		Password: "ali",
	}
	user3 := entity.User{
		Username: "",
		Password: "ali",
	}
	user4 := entity.User{
		Username: "mamad",
		Password: "",
	}

	_, err = authenticationService.Register(user1.Username, user1.Password)
	if err != nil {
		t.Error(err)
	}

	_, err = authenticationService.Register(user2.Username, user2.Password)
	if err != nil && !errors.Is(err, messages.ErrAlreadyExists) {
		t.Error(err)
	}
	if err == nil {
		t.Error(err)
	}

	_, err = authenticationService.Register(user3.Username, user3.Password)
	if err != nil && !errors.Is(err, messages.ErrBadRequest) {
		t.Error(err)
	}

	_, err = authenticationService.Register(user4.Username, user4.Password)
	if err != nil && !errors.Is(err, messages.ErrBadRequest) {
		t.Error(err)
	}
}
