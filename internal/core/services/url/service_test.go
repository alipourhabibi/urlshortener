package url_test

import (
	"testing"

	"github.com/alipourhabibi/urlshortener/config"
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/alipourhabibi/urlshortener/internal/core/services/url"
	"github.com/google/uuid"
)

func TestURL(t *testing.T) {
	err := config.Confs.Load("./local_config.yaml")
	if err != nil {
		t.Error(err)
	}
	urlService, err := url.New(url.WithPostgresURLRepository(config.Confs.PostgresDB), url.WithRedisCacheRepository(config.Confs.Redis))
	if err != nil {
		t.Error(err)
	}
	mockUser := entity.User{
		ID:       uuid.New(),
		Username: "ali",
		Password: "ali",
	}
	mockURLS := []entity.Shortened{
		{
			OriginalURL: "https://google.com",
		},
		{
			OriginalURL: "https://yahoo.com",
		},
		{
			OriginalURL: "https://alibaba.com",
		},
	}
	for k := range mockURLS {
		shortened, err := urlService.Add(mockURLS[k].OriginalURL, &mockUser)
		if err != nil {
			t.Error(err)
		}
		mockURLS[k].ShortenedURL = shortened
	}
	for k := range mockURLS {
		url, err := urlService.Get(mockURLS[k].ShortenedURL)
		if err != nil {
			t.Error(err)
		}
		if mockURLS[k].OriginalURL != url.GetOriginal() {
			t.Errorf("Mishmatched original and shortened")
		}
	}
}
