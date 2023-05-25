package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nanwp/travello/config"
	"github.com/nanwp/travello/helper"
	"github.com/nanwp/travello/models/users"
	"github.com/nanwp/travello/pkg/middleware/auth"
	"github.com/nanwp/travello/service"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) VerifyEmail(c *gin.Context) {
	userId := c.Query("id")

	if userId == "" {
		c.HTML(http.StatusNotFound, "not-found.html", gin.H{})
		return
	}

	user, err := h.userService.VerifyEmail(userId)

	if err != nil {
		c.HTML(http.StatusOK, "not-found.html", gin.H{
			"data": user,
		})
		return
	}

	c.HTML(http.StatusOK, "verification-success.html", gin.H{
		"data": user,
	})
}

func (h *userHandler) Register(c *gin.Context) {
	var userRequest users.UserCreate

	err := c.ShouldBindJSON(&userRequest)
	if err != nil {
		errorMessages := []string{}

		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error on fieled %s, conditions: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", errorMessages, nil)
		return
	}

	emailCheck, _ := h.userService.FindAll()

	for _, a := range emailCheck {
		if a.Email == userRequest.Email {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", "email telah digunakan", nil)
			return
		}
	}

	user, err := h.userService.Create(userRequest)

	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err, nil)
		return
	}

	helper.ResponseOutput(c, http.StatusCreated, "CREATED", "check your email to be verified", user)
}

func (h *userHandler) GetUser(c *gin.Context) {

	id := auth.UserID

	user, err := h.userService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err,
			"status":  "BAD_REQUEST",
		})
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err, nil)
		return
	}

	helper.ResponseOutput(c, http.StatusOK, "OK", nil, user)

}

func (h *userHandler) UpdateUser(c *gin.Context) {
	id := auth.UserID
	var userUpdate users.UserUpdate

	err := c.ShouldBindJSON(&userUpdate)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err, nil)
		return
	}

	user, err := h.userService.Update(id, userUpdate)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err, nil)
		return
	}

	responseData := users.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		Password: user.Password,
	}
	helper.ResponseOutput(c, http.StatusOK, "OK", "success update data", responseData)
}

func (h *userHandler) UpdatePassword(c *gin.Context) {
	id := auth.UserID
	var passwordUpdate users.UpdatePassword
	err := c.ShouldBindJSON(&passwordUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err,
		})
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err, nil)
		return
	}

	user, err := h.userService.FindByID(id)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordUpdate.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": "old password not match",
		})
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", "old password not match", nil)
		return
	}

	updated, err := h.userService.UpdatePassword(auth.UserID, passwordUpdate.NewPassword)

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "success update password",
		"data":    updated,
	})
	helper.ResponseOutput(c, http.StatusOK, "OK", "success update password", updated)

}

func (h *userHandler) Login(c *gin.Context) {
	var userInput users.UserLogin

	err := c.ShouldBindJSON(&userInput)
	if err != nil {
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("error on field %s, conditions: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", errorMessages, nil)
		return
	}

	userLogin, err := h.userService.FindByEmail(userInput.Email)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case gorm.ErrRecordNotFound:
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", "user not found", nil)
			return

		default:
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err, nil)

			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userLogin.Password), []byte(userInput.Password)); err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", "password not match", nil)
		return
	}

	if userLogin.Verified {
		expTime := time.Now().Add(time.Hour * 8640)
		claims := &config.JWTClaim{
			UserID: userLogin.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "trevello",
				ExpiresAt: jwt.NewNumericDate(expTime),
			},
		}

		tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		token, err := tokenAlgo.SignedString(config.JWT_KEY)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "INTERNAL_SERVER_ERROR",
				"message": err.Error(),
			})
			helper.ResponseOutput(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err, nil)
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "token",
			Value:    token,
			MaxAge:   3600 * 24 * 30,
			Path:     "/",
			HttpOnly: true,
		})

		type resp struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Token string `json:"token"`
		}

		datRes := resp{
			Name:  userLogin.Name,
			Email: userLogin.Email,
			Token: token,
		}

		helper.ResponseOutput(c, http.StatusOK, "OK", "sucess login", datRes)
		return
	}

	helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", "user not verified, please check your email", nil)

}
