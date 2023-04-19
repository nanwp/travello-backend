package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nanwp/travello/config"
	"github.com/nanwp/travello/middleware"
	"github.com/nanwp/travello/models/users"
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

func (h *userHandler) Register(c *gin.Context) {
	var userRequest users.UserCreate

	err := c.ShouldBindJSON(&userRequest)
	if err != nil {
		errorMessages := []string{}

		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error on fieled %s, conditions: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": errorMessages,
		})
		return
	}

	emailCheck, _ := h.userService.FindAll()

	for _, a := range emailCheck {
		if a.Email == userRequest.Email {
			c.JSON(http.StatusBadRequest, gin.H{
				"email":   a.Email,
				"status":  "BAD_REQUEST",
				"message": "email telah digunakan",
			})
			return
		}
	}

	user, err := h.userService.Create(userRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
			"status":  "BAD_REQUEST",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "succes create user",
		"data":    user,
	})
}

func (h *userHandler) GetUser(c *gin.Context) {

	id := middleware.UserID

	user, err := h.userService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err,
			"status":  "BAD_REQUEST",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "success get data",
		"data":    user,
	})

}

func (h *userHandler) UpdateUser(c *gin.Context) {
	id := middleware.UserID
	var userUpdate users.UserUpdate

	err := c.ShouldBindJSON(&userUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err,
		})
		return
	}

	user, err := h.userService.Update(id, userUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err,
		})
		return
	}

	responseData := users.UserResponse{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		Role: user.Role,
		Password: user.Password,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"message":"success update data",
		"data":   responseData,
	})
}

func (h *userHandler) UpdatePassword(c *gin.Context) {
	id := middleware.UserID
	var passwordUpdate users.UpdatePassword
	err := c.ShouldBindJSON(&passwordUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err,
		})
		return
	}

	user, err := h.userService.FindByID(id)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordUpdate.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": "old password not match",
		})
		return
	}



	updated, err := h.userService.UpdatePassword(middleware.UserID, passwordUpdate.NewPassword)


	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"message":"success update password",
		"data":   updated,
	})

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

		c.JSON(http.StatusBadRequest, gin.H{
			"status": "BAD_REQUEST",
			"errors": errorMessages,
		})
		return
	}

	userLogin, err := h.userService.FindByEmail(userInput.Email)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "BAD_REQUEST",
				"message": "user not found",
			})
			return

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "BAD_REQUEST",
				"message": err,
			})
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userLogin.Password), []byte(userInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "UNAUTHORIZED",
			"message": "password not match",
		})
		return
	}

	expTime := time.Now().Add(time.Hour * 24)
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
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "sucess login",
		"token":   token,
	})
}
