package service

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/nanwp/travello/config"
	"github.com/nanwp/travello/models/users"
	"github.com/nanwp/travello/repository"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type UserService interface {
	Create(user users.UserCreate) (users.User, error)
	FindByEmail(email string) (users.User, error)
	FindAll() ([]users.User, error)
	FindByID(id string) (users.UserResponse, error)
	Update(ID string, userUpdate users.UserUpdate) (users.User, error)
	UpdatePassword(ID string, newPassword string) (users.User, error)
	VerifyEmail(ID string) (users.UserResponse, error)
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
		Name:     strings.Title(strings.ToLower(fullName)),
		Email:    user.Email,
		Password: string(hashPassword),
		Role:     "user",
		Verified: false,
	}

	newUser, err := s.repository.Create(usr)

	go sendVerificationEmail(usr)

	return newUser, err
}

func (s *userService) FindByEmail(email string) (users.User, error) {
	return s.repository.FindByEmail(email)
}

func (s *userService) FindByID(id string) (users.UserResponse, error) {
	userGet, err := s.repository.FindByID(id)

	user := users.UserResponse{
		ID:       userGet.ID,
		Name:     userGet.Name,
		Email:    userGet.Email,
		Role:     userGet.Role,
		Password: userGet.Password,
	}

	return user, err
}

func (s *userService) FindAll() ([]users.User, error) {
	return s.repository.FindAll()
}

func (s *userService) VerifyEmail(ID string) (users.UserResponse, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return users.UserResponse{}, err
	}

	if user.Verified {
		err := errors.New("user telah terverifikasi")
		resp := users.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			Role:     user.Role,
			Verified: user.Verified,
		}
		return resp, err
	}

	user.Verified = true

	updateUser, err := s.repository.Update(user)

	resp := users.UserResponse{
		ID:       updateUser.ID,
		Name:     updateUser.Name,
		Email:    updateUser.Email,
		Password: updateUser.Password,
		Role:     updateUser.Role,
		Verified: updateUser.Verified,
	}
	return resp, nil

}

func sendVerificationEmail(user users.User) error {
	url := fmt.Sprintf("http://103.171.182.206:8070/verify?id=%s", user.ID)
	bodyEmail := fmt.Sprintf("Hello %s,<br><br>Klik link di bawah ini untuk mengaktifkan akun kamu:<br><br><a href=\"%s\">Aktifkan Akun</a>", user.Name, url)
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.CONFIG_SENDER_NAME)
	mailer.SetHeader("To", user.Email)
	mailer.SetHeader("Subject", "Verif email")
	mailer.SetBody("text/html", bodyEmail)

	dialer := gomail.NewDialer(
		config.CONFIG_SMTP_HOST,
		config.CONFIG_SMTP_PORT,
		config.CONFIG_AUTH_EMAIL,
		config.CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Printf("error send email : %s", err.Error())
	}

	log.Println("Success send!")
	return nil
}
