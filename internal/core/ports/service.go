package ports

import (
	aggregate "github.com/alipourhabibi/urlshortener/internal/core/aggergate"
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/google/uuid"
)

type URLService interface {
	Add(string, *entity.User) error
	Get(uuid.UUID) (aggregate.URL, error)
}
