package database

import (
	"log"
	"os"

	"github.com/Brawdunoir/dionysos-server/models"
	c "github.com/Brawdunoir/dionysos-server/variables"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB = setupDatabase()

// setupDatabase returns a setuped database instance.
func setupDatabase() *gorm.DB {
	resetDB := c.Environment == c.ENVIRONMENT_DEVELOPMENT || c.Environment == c.ENVIRONMENT_TESTING

	db, err := gorm.Open(postgres.Open(createDSN()), createConfig())
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	err = MigrateDB(db, resetDB)
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

// createConfig creates a Gorm config depending on the environment variables.
func createConfig() *gorm.Config {
	switch c.Environment {
	case c.ENVIRONMENT_TESTING:
		return &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
	case c.ENVIRONMENT_DEVELOPMENT:
		return &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	case c.ENVIRONMENT_PRODUCTION:
		return &gorm.Config{Logger: logger.Default.LogMode(logger.Error)}
	default:
		log.Println("ENVIRONMENT variable not valid, using default config")
		log.Println("Possible values are : " + c.ENVIRONMENT_TESTING + ", " + c.ENVIRONMENT_DEVELOPMENT + ", " + c.ENVIRONMENT_PRODUCTION)
		return &gorm.Config{}
	}
}

// MigrateDB migrate a table in the database and resets all tables if needed.
func MigrateDB(db *gorm.DB, reset bool) error {
	if reset {
		err := db.Migrator().DropTable(&models.User{}, &models.Room{}, "room_users")
		if err != nil {
			return err
		}
	}
	err := db.AutoMigrate(&models.Room{}, &models.User{})
	if err != nil {
		return err
	}
	return nil
}
