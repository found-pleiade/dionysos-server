package main

import (
	"os"

	"github.com/Brawdunoir/dionysos-server/docs"
	"github.com/Brawdunoir/dionysos-server/routes"
	"github.com/Brawdunoir/dionysos-server/utils"
	"github.com/Brawdunoir/dionysos-server/variables"
	"github.com/gin-gonic/gin"
)

// VERSION as a constant variable
const VERSION = "0.3.0"

// @title           Dionysos
// @description     API instance for the Dionysos client application.
// @securityDefinitions.basic BasicAuth

// @contact.name   API Support
// @contact.url    https://github.com/Brawdunoir/dionysos-server/issues

// @license.name  GNU General Public License v3.0
// @license.url   https://www.gnu.org/licenses/gpl-3.0.html
func main() {
	os.Setenv("VERSION", VERSION)

	// Set dynamic documentation information.
	docs.SwaggerInfo.Version = VERSION
	docs.SwaggerInfo.BasePath = variables.BasePath

	// Logger initialization.
	err := utils.InitLogger()
	if err != nil {
		panic(err)
	}
	defer utils.Logger.Sync()

	// Gin initialization.
	router := routes.SetupRouter(gin.New())

	// Start the server.
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
