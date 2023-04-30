package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/jub0bs/fcors"
	"github.com/nanwp/travello/config"
	"github.com/nanwp/travello/handler"
	"github.com/nanwp/travello/middleware"
	"github.com/nanwp/travello/repository"
	"github.com/nanwp/travello/service"
)

func main() {

	r := gin.Default()
	db := config.ConnectDatabase()

	userHandler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(db)))
	destinatinHandler := handler.NewDestinationHandler()
	ulasanHandler := handler.NewUlasanHandler()

	cors, err := fcors.AllowAccess(
		fcors.FromAnyOrigin(),
		fcors.WithMethods(
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			"UPDATE",
		),
		fcors.WithRequestHeaders(
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Max",
		),
		fcors.MaxAgeInSeconds(86400),
	)
	if err != nil {
		log.Fatal(err)
	}
	r.Use(adapter.Wrap(cors))

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.GET("/destination", destinatinHandler.Destinations)
	r.GET("/user", middleware.JWTMiddleware, userHandler.GetUser)
	r.PUT("/user", middleware.JWTMiddleware, userHandler.UpdateUser)
	r.PUT("/userpassword", middleware.JWTMiddleware, userHandler.UpdatePassword)

	r.POST("/destination", destinatinHandler.Create)
	// r.GET("/destination", destinatinHandler.DestinationCategory)

	r.POST("/ulasan", middleware.JWTMiddleware, ulasanHandler.AddUlasan)
	r.GET("/ulasan", middleware.JWTMiddleware, ulasanHandler.GetUlasanByDestination)

	r.Run(":8080")

}
