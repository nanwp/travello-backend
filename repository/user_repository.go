package repository

import (
	"github.com/nanwp/travello/models/users"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user users.User) (users.User, error)
	FindByEmail(email string) (users.User, error)
	FindAll() ([]users.User, error)
	FindByID(id string) (users.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user users.User) (users.User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

func (r *userRepository) FindByEmail(email string) (users.User, error) {
	var user users.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) FindByID(id string) (users.User, error) {
	var user users.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *userRepository) FindAll() ([]users.User, error) {
	var users []users.User
	err := r.db.Find(&users).Error
	return users, err
}
