package repository

import (
	"github.com/nanwp/travello/models/favorites"
	"gorm.io/gorm"
)

type FavoriteRepository interface {
	Create(favorite favorites.Favorite) error
	CheckFavorit(destinationId string, userId string) (bool, error)
	Delete(favorite favorites.Favorite) error
	GetFavoriteByUser(userId string) ([]favorites.Favorite, error)
}

type favoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *favoriteRepository {
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) Create(favorite favorites.Favorite) error {
	return r.db.Create(&favorite).Error
}

func (r *favoriteRepository) CheckFavorit(destinationId string, userId string) (bool, error) {
	model := &favorites.Favorite{}
	var exists bool
	err := r.db.Model(model).Select("count(*) > 0").Where("destination_id = ? AND user_id = ?", destinationId, userId).Find(&exists).Error
	return exists, err
}

func (r *favoriteRepository) Delete(favorite favorites.Favorite) error {
	return r.db.Delete(favorite).Error
}

func (r *favoriteRepository) GetFavoriteByUser(userId string) ([]favorites.Favorite, error) {
	var favorites []favorites.Favorite
	err := r.db.Where("user_id = ?", userId).Find(&favorites).Error
	return favorites, err
}
