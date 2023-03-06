package postgres

import (
	"fmt"

	"github.com/alipourhabibi/urlshortener/config"
	aggregate "github.com/alipourhabibi/urlshortener/internal/core/aggergate"
	urldomain "github.com/alipourhabibi/urlshortener/internal/core/domain/url"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func New(confs config.Posgres) (*PostgresRepository, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		confs.Host, confs.User, confs.Password,
		confs.DBName, confs.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return &PostgresRepository{
		db: db,
	}, err
}

func (p *PostgresRepository) Add(u aggregate.URL) error {
	return p.db.Model(urldomain.PostgresURL{}).Create(u).Error
}

func (p *PostgresRepository) Get(uuid uuid.UUID) (aggregate.URL, error) {
	pu := urldomain.PostgresURL{}
	err := p.db.Model(urldomain.PostgresURL{}).Find(&pu, "id = ?", uuid).Error
	if err != nil {
		return aggregate.URL{}, err
	}
	url := pu.ToAggregate()
	return url, nil
}
