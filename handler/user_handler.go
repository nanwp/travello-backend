package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nanwp/travello/models/users"
	"github.com/nanwp/travello/service"
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
			"error": errorMessages,
		})
		return
	}

	emailCheck, _ := h.userService.FindAll()

	for _, a := range emailCheck {
		if a.Email == userRequest.Email {
			c.JSON(http.StatusBadRequest, gin.H{
				"email": a.Email,
				"error": "email telah digunakan",
			})
			return
		}
	}

	user, err := h.userService.Create(userRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
