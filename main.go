package main

import (
	"os"

	"github.com/Brawdunoir/dionysos-server/routes"
	"github.com/gin-gonic/gin"
)

// VERSION as a constant variable
const VERSION = "0.2.0"

func main() {
	os.Setenv("VERSION", VERSION)

	router := routes.SetupRouter(gin.Default())

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
