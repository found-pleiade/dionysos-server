package utils

import (
	"log"

	"github.com/Brawdunoir/dionysos-server/database"
	l "github.com/Brawdunoir/dionysos-server/utils/logger"
	"github.com/Brawdunoir/dionysos-server/variables"
)

func InitAPI() {
	// ! Order of these calls matters.
	// Load variables from environment.
	variables.LoadVariables()

	// Logger initialization.
	err := l.InitLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Logger.Sync()
	l.Logger.Debug("Logger initialized.")

	// DB initialization.
	database.Init()
	l.Logger.Debug("Database initialized.")
}
