package main

import (
	"os"

	"github.com/Brawdunoir/dionysos-server/routes"
)

// VERSION as a constant variable
const VERSION = "0.1.0"

func main() {
	os.Setenv("VERSION", VERSION)

	router := routes.SetupRouter()

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
