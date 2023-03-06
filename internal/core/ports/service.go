package ports

import (
	aggregate "github.com/alipourhabibi/urlshortener/internal/core/aggergate"
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
)

type URLService interface {
	Add(string, *entity.User) (string, error)
	Get(string) (aggregate.URL, error)
}
