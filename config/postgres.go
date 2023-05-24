package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	url := "postgres://postgres:Latihan@103.171.182.206:5432/travello"
	database, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	// dbInstance, err := database.DB()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer dbInstance.Close()

	return database
}
