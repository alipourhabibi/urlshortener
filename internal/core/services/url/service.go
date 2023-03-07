package url

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/alipourhabibi/urlshortener/config"
	aggregate "github.com/alipourhabibi/urlshortener/internal/core/aggergate"
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/alipourhabibi/urlshortener/internal/core/ports"
	postgresdb "github.com/alipourhabibi/urlshortener/internal/repository/postgres"
	redisdb "github.com/alipourhabibi/urlshortener/internal/repository/redis"
	"github.com/go-redis/redis"
)

type UrlConfiguration func(*UrlService) error

type UrlService struct {
	urlRepository   ports.URLRepository
	cacheRepository ports.CacheRepository
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

func WithRedisCacheRepository(config config.Redis) UrlConfiguration {
	return func(us *UrlService) error {
		client, err := redisdb.New(config)
		if err != nil {
			return err
		}
		us.cacheRepository = client
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
	_, err = s.cacheRepository.Set(shortened, original, 2*time.Minute)
	if err != nil {
		return "", err
	}
	return s.urlRepository.Add(url)
}

func (s *UrlService) Get(shortened string) (aggregate.URL, error) {
	original, err := s.cacheRepository.Get(shortened)
	if err != nil && !errors.Is(err, redis.Nil) {
		return aggregate.URL{}, err
	}
	if !errors.Is(err, redis.Nil) {
		au := aggregate.NewURLEmpty()
		au.SetOriginal(original.(string))
		au.SetShortened(shortened)
		return *au, nil
	}
	url, err := s.urlRepository.Get(shortened)
	if err != nil {
		return aggregate.URL{}, err
	}
	_, err = s.cacheRepository.Set(shortened, url.GetOriginal(), 2*time.Minute)
	if err != nil {
		return aggregate.URL{}, err
	}
	return url, nil
}
