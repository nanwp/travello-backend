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
	FindByID(id string) (users.UserResponse, error)
	Update(ID string, userUpdate users.UserUpdate) (users.User, error)
	UpdatePassword(ID string, newPassword string) (users.User, error)
}

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *userService {
	return &userService{repository}
}
func (s *userService) UpdatePassword(ID string, newPassword string) (users.User, error) {

	user, err := s.repository.FindByID(ID)

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	user.Password = string(hashPassword)

	updated, err := s.repository.Update(user)

	return updated, err
}

func (s *userService) Update(ID string, userUpdate users.UserUpdate) (users.User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	if userUpdate.Name != "" {
		user.Name = userUpdate.Name
	}

	if userUpdate.Email != "" {
		user.Email = userUpdate.Email
	}

	updateUser, err := s.repository.Update(user)

	return updateUser, err
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

func (s *userService) FindByID(id string) (users.UserResponse, error) {
	userGet, err := s.repository.FindByID(id)

	user := users.UserResponse{
		ID:    userGet.ID,
		Name:  userGet.Name,
		Email: userGet.Email,
		Role:  userGet.Role,
		Password: userGet.Password,
	}

	return user, err
}

func (s *userService) FindAll() ([]users.User, error) {
	return s.repository.FindAll()
}
