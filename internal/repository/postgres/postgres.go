package postgres

import (
	"fmt"

	"github.com/alipourhabibi/urlshortener/config"
	aggregate "github.com/alipourhabibi/urlshortener/internal/core/aggergate"
	urldomain "github.com/alipourhabibi/urlshortener/internal/core/domain/url"
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

func (p *PostgresRepository) Add(u aggregate.URL) (string, error) {
	payload := urldomain.PostgresURL{}
	payload.ID = u.GetID()
	payload.Original = u.GetOriginal()
	payload.Shortened = u.GetShortened()
	err := p.db.Model(urldomain.PostgresURL{}).Create(&payload).Error
	if err != nil {
		return "", err
	}
	return payload.Shortened, nil
}

func (p *PostgresRepository) Get(shortened string) (aggregate.URL, error) {
	pu := urldomain.PostgresURL{}
	err := p.db.Model(urldomain.PostgresURL{}).Find(&pu, "shortened = ?", shortened).Error
	if err != nil {
		return aggregate.URL{}, err
	}
	url := pu.ToAggregate()
	return url, nil
}

func (p *PostgresRepository) Exists(shotened string) (bool, error) {
	var exists bool
	err := p.db.Model(urldomain.PostgresURL{}).Select("count(*) > 0").Where("shortened = ?", shotened).Find(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p *PostgresRepository) MigrateAll() error {
	fmt.Printf("Migrating PostgresURL\n")
	err := p.db.AutoMigrate(&urldomain.PostgresURL{})
	if err != nil {
		return err
	}
	return nil
}
