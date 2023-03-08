package postgres

import (
	"errors"
	"fmt"

	"github.com/alipourhabibi/urlshortener/config"
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/alipourhabibi/urlshortener/internal/core/messages"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(confs config.Posgres) (*AuthRepository, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		confs.Host, confs.User, confs.Password,
		confs.DBName, confs.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return &AuthRepository{
		db: db,
	}, err
}

func (a *AuthRepository) Add(u entity.User) error {
	return a.db.Model(entity.User{}).Create(&u).Error
}

func (a *AuthRepository) Exists(username string) (bool, error) {
	var exists bool
	err := a.db.Model(entity.User{}).Select("count(*) > 0").Where("username = ?", username).Find(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (a *AuthRepository) Get(username string) (entity.User, error) {
	u := entity.User{}
	err := a.db.Model(entity.User{}).Find(&u, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, messages.ErrNotFound
		}
		return entity.User{}, messages.ErrInternalServerError
	}
	return u, nil
}
