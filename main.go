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
	destinatinHandler := handler.NewDestinationHandler()

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.GET("/destination", destinatinHandler.Destinations)
	r.POST("/destination", destinatinHandler.Create)
	// r.GET("/destination", destinatinHandler.DestinationCategory)

	r.Run(":8080")

}
