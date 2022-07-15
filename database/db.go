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
	db, err := gorm.Open(postgres.Open(createDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// nolint:errcheck
	db.AutoMigrate(&models.User{}, &models.Room{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	return db
}

// createDSN creates a DSN string from the environment variables to connect to the database.
func createDSN() string {
	username, found := os.LookupEnv("POSTGRES_USER")
	if !found {
		log.Fatal("POSTGRES_USER environment variable not found")
	}
	password, found := os.LookupEnv("POSTGRES_PASSWORD")
	if !found {
		log.Fatal("POSTGRES_PASSWORD environment variable not found")
	}
	host, found := os.LookupEnv("POSTGRES_HOST")
	if !found {
		log.Fatal("POSTGRES_HOST environment variable not found")
	}
	port, found := os.LookupEnv("POSTGRES_PORT")
	if !found {
		log.Fatal("POSTGRES_PORT environment variable not found")
	}
	dbname, found := os.LookupEnv("POSTGRES_DB")
	if !found {
		log.Fatal("POSTGRES_DB environment variable not found")
	}

	return "host=" + host + " port=" + port + " user=" + username + " password=" + password + " dbname=" + dbname
}
