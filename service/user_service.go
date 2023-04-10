package service

import (
	"github.com/google/uuid"
	"github.com/nanwp/travello/models/users"
	"github.com/nanwp/travello/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(user users.UserCreate) (users.User, error)
	FindByEmail(email string) (users.User, error)
	FindAll() ([]users.User, error)
}

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *userService {
	return &userService{repository}
}

func (s *userService) Create(user users.UserCreate) (users.User, error) {
	uuidGenerate := uuid.New()
	fullName := user.FirstName + " " + user.LastName
	stringUuid := uuidGenerate.String()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	usr := users.User{
		ID:       stringUuid,
		Name:     fullName,
		Email:    user.Email,
		Password: string(hashPassword),
		Role:     "user",
	}

	newUser, err := s.repository.Create(usr)
	return newUser, err
}

func (s *userService) FindByEmail(email string) (users.User, error) {
	return s.repository.FindByEmail(email)
}

func (s *userService) FindAll() ([]users.User, error) {
	return s.repository.FindAll()
}
