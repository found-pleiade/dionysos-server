package database

import (
	"github.com/Brawdunoir/dionysos-server/models"
	l "github.com/Brawdunoir/dionysos-server/utils/logger"
	c "github.com/Brawdunoir/dionysos-server/variables"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var database *gorm.DB

// Init initializes the database connection and migrate models within the database.
// It also resets the database in case of testing or development environment.
func Init() {
	resetDB := c.Environment == c.ENVIRONMENT_DEVELOPMENT || c.Environment == c.ENVIRONMENT_TESTING

	db, err := gorm.Open(postgres.Open(createDSN()), createConfig())
	if err != nil {
		l.Logger.Fatal("Failed to connect to the database: ", err)
	}

	err = MigrateDB(db, resetDB)
	if err != nil {
		l.Logger.Fatal("Failed to migrate database: ", err)
	}
	database = db
}

// GetDB returns the database connection.
func GetDB() *gorm.DB {
	return database
}

// createDSN creates a DSN string from the environment variables to connect to the database.
func createDSN() string {
	return "host=" + c.PostgresHost + " port=" + c.PostgresPort + " user=" + c.PostgresUser + " password=" + c.PostgresPassword + " dbname=" + c.PostgresDB
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
		l.Logger.Fatal("Unknown environment: " + c.Environment)
		return nil
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
