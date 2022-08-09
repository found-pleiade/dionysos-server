package routes

import (
	"github.com/Brawdunoir/dionysos-server/database"
	"gorm.io/gorm"
)

// db is the database used by the api across all routes.
var db *gorm.DB = database.DB //nolint:deadcode,unused,varcheck
