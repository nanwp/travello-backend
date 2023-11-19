package main

import (
	"fmt"

	"github.com/nanwp/travello/config"
	"github.com/nanwp/travello/pkg/middleware"
)

func main() {
	db := config.ConnectDatabase()

	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	router := middleware.InitRouter(db)

	if err := router.Run(":8180"); err != nil {
		panic(fmt.Errorf("failed to start server: %s", err))
	}
}
