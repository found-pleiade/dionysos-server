package routes

import (
	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// db is the database used by the api accross all routes.
var db *gorm.DB = database.GetDatabase() //nolint:deadcode,unused,varcheck

// validate is the validator used by the api accross all routes.
var validate *validator.Validate = validator.New() //nolint:deadcode,unused,varcheck
