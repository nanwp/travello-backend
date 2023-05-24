package repository

import (
	"github.com/nanwp/travello/models/ulasans"
	"gorm.io/gorm"
)

type UlasanRepository interface {
	Create(ulasan ulasans.Ulasan) (ulasans.Ulasan, error)
	GetByDestinationID(id string) ([]ulasans.Ulasan, error)
	GetAllUlasan() ([]ulasans.Ulasan, error)
}

type ulasanRepository struct {
	db *gorm.DB
}

func NewUlasanRepository(db *gorm.DB) *ulasanRepository {
	return &ulasanRepository{db}
}

func (r *ulasanRepository) Create(ulasan ulasans.Ulasan) (ulasans.Ulasan, error) {
	err := r.db.Create(&ulasan).Error
	return ulasan, err
}

func (r *ulasanRepository) GetAllUlasan() ([]ulasans.Ulasan, error) {
	var ulasan []ulasans.Ulasan
	err := r.db.Preload("User").Find(&ulasan).Error
	return ulasan, err
}

func (r *ulasanRepository) GetByDestinationID(id string) ([]ulasans.Ulasan, error) {
	var ulasan []ulasans.Ulasan
	err := r.db.Where("destination_id = ?", id).Preload("User").Find(&ulasan).Error
	return ulasan, err
}
