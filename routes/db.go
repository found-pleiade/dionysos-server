package routes

import (
	"github.com/Brawdunoir/dionysos-server/database"
	"gorm.io/gorm"
)

// db is the database used by the api accross all routes.
var db *gorm.DB = database.GetDatabase() //nolint:deadcode,unused,varcheck
