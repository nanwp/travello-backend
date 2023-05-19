package main

import (
	"fmt"

	"github.com/nanwp/travello/config"
	"github.com/nanwp/travello/pkg/middleware"
)

func main() {

	db := config.ConnectDatabase()

	router := middleware.InitRouter(db)

	if err := router.Run(":8080"); err != nil {
		panic(fmt.Errorf("failed start server: %s", err))
	}
}
