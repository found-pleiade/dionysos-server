package routes

import (
	"github.com/Brawdunoir/dionysos-server/database"

	"github.com/arangodb/go-driver"
)

// db is the database used by the api accross all routes
var db driver.Database = database.GetDatabase("dionysos")
