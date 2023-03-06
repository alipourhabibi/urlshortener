package url

import (
	"github.com/alipourhabibi/urlshortener/config"
	aggregate "github.com/alipourhabibi/urlshortener/internal/core/aggergate"
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/alipourhabibi/urlshortener/internal/core/ports"
	repo "github.com/alipourhabibi/urlshortener/internal/repository"
	"github.com/google/uuid"
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
		client, err := repo.New(config)
		if err != nil {
			return err
		}
		us.urlRepository = client
		return nil
	}
}

func (s *UrlService) GenerateShoretend(original string) (string, error) {
	return "", nil
}

func (s *UrlService) Add(original string, user *entity.User) error {
	shortened, err := s.GenerateShoretend(original)
	if err != nil {
		return err
	}

	url, err := aggregate.NewURL(original, user)
	if err != nil {
		return err
	}
	url.SetShortened(shortened)
	return s.urlRepository.Add(url)
}

func (s *UrlService) Get(id uuid.UUID) (aggregate.URL, error) {
	return s.urlRepository.Get(id)
}
