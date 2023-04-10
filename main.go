package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nanwp/travello/config"
	"github.com/nanwp/travello/handler"
	"github.com/nanwp/travello/repository"
	"github.com/nanwp/travello/service"
)

func main() {

	r := gin.Default()
	db := config.ConnectDatabase()

	userHandler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(db)))

	r.POST("/daftar", userHandler.Register)

	r.Run(":8080")
}
