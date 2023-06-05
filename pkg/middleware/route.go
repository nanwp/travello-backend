package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nanwp/travello/handler"
	"github.com/nanwp/travello/pkg/middleware/auth"
	"github.com/nanwp/travello/repository"
	"github.com/nanwp/travello/service"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	router := gin.New()
	router.Use(CORSMiddleware())

	userRepo := repository.NewUserRepository(db)
	ulasanRepo := repository.NewUlasanRepository(db)

	userService := service.NewUserService(userRepo)
	ulasanService := service.NewUlasanService(ulasanRepo, *userService)

	userHandler := handler.NewUserHandler(userService)
	destinatinHandler := handler.NewDestinationHandler(ulasanService)
	ulasanHandler := handler.NewUlasanHandler(ulasanService)

	r := router.Group("api")

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	// r.GET("/verify", userHandler.VerifyEmail)

	r.GET("/destinations", destinatinHandler.Destinations)
	r.POST("/destination", destinatinHandler.Create)
	r.GET("/destination/:id", destinatinHandler.Destination)

	r.GET("/user", auth.JWTMiddleware, userHandler.GetUser)
	r.PUT("/user", auth.JWTMiddleware, userHandler.UpdateUser)
	r.PUT("/userpassword", auth.JWTMiddleware, userHandler.UpdatePassword)

	r.POST("/ulasan", auth.JWTMiddleware, ulasanHandler.AddUlasan)
	// r.GET("/ulasan", ulasanHandler.GetUlasanByDestination)

	router.GET("/verify", userHandler.VerifyEmail)
	router.LoadHTMLGlob("templates/*/*.html")

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "Page not found",
		})
	})

	return router

}
