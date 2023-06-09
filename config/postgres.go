package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const usernamePostgre = "postgres"
const passwordPostgre = "NewVPSNanda"

func ConnectDatabase() *gorm.DB {
	url := "postgres://" + usernamePostgre + ":" + passwordPostgre + "@103.161.184.72:5432/travello"
	database, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	return database
}
