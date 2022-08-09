package routes

import (
	"github.com/Brawdunoir/dionysos-server/database"
	"gorm.io/gorm"
)

// db is the database used by the api across all routes.
var db *gorm.DB = database.GetDatabase() //nolint:deadcode,unused,varcheck

// GetDB returns the database instance.
func GetDB() *gorm.DB {
	return db
}
