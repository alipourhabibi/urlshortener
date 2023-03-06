package url

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"

	"github.com/alipourhabibi/urlshortener/config"
	aggregate "github.com/alipourhabibi/urlshortener/internal/core/aggergate"
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/alipourhabibi/urlshortener/internal/core/ports"
	postgresdb "github.com/alipourhabibi/urlshortener/internal/repository/postgres"
)

type UrlConfiguration func(*UrlService) error

type UrlService struct {
	urlRepository ports.URLRepository
}

func New(cfgs ...UrlConfiguration) (*UrlService, error) {
	us := &UrlService{}

	// Loop through all the cfgs and apply them
	for _, cfg := range cfgs {
		err := cfg(us)
		if err != nil {
			return nil, err
		}
	}
	return us, nil
}

func WithPostgresURLRepository(config config.Posgres) UrlConfiguration {
	return func(us *UrlService) error {
		client, err := postgresdb.New(config)
		if err != nil {
			return err
		}
		us.urlRepository = client
		return nil
	}
}

func (s *UrlService) generateShoretend(original string, i, j uint) (string, error) {
	sum := md5.Sum([]byte(original))
	sumString := fmt.Sprintf("%s", sum)
	encoded := base64.StdEncoding.EncodeToString([]byte(sumString))

	return encoded[i:j], nil
}

func (s *UrlService) Add(original string, user *entity.User) (string, error) {
	i, j := uint(0), uint(6)
	shortened, err := s.generateShoretend(original, i, j)
	if err != nil {
		return "", err
	}
	exists, err := s.urlRepository.Exists(shortened)
	if err != nil {
		return "", err
	}
	for exists {
		i, j = i+2, j+2
		shortened, err = s.generateShoretend(original, i, j)
		if err != nil {
			return "", err
		}
		exists, err = s.urlRepository.Exists(shortened)
		if err != nil {
			return "", err
		}
	}

	url, err := aggregate.NewURL(original, user)
	if err != nil {
		return "", err
	}
	url.SetShortened(shortened)
	return s.urlRepository.Add(url)
}

func (s *UrlService) Get(shortened string) (aggregate.URL, error) {
	return s.urlRepository.Get(shortened)
}
