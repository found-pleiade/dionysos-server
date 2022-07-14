package database

import (
	"log"
	"os"

	"github.com/Brawdunoir/dionysos-server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetDatabase returns a database instance
func GetDatabase() *gorm.DB {
	username, found := os.LookupEnv("POSTGRES_USER")
	if !found {
		log.Fatal("POSTGRES_USER environment variable not found")
	}
	password, found := os.LookupEnv("POSTGRES_PASSWORD")
	if !found {
		log.Fatal("POSTGRES_PASSWORD environment variable not found")
	}

	dsn := "host=postgres port=5432 user=" + username + " password=" + password + " dbname=dionysos"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// nolint:errcheck
	db.AutoMigrate(&models.User{}, &models.Room{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	return db
}
