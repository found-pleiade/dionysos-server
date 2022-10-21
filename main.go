package main

import (
	"os"

	"github.com/Brawdunoir/dionysos-server/docs"
	"github.com/Brawdunoir/dionysos-server/routes"
	"github.com/Brawdunoir/dionysos-server/utils"
	l "github.com/Brawdunoir/dionysos-server/utils/logger"
	"github.com/Brawdunoir/dionysos-server/variables"
	"github.com/gin-gonic/gin"
)

// VERSION as a constant variable
const VERSION = "0.4.0"

// @title           Dionysos
// @description     API instance for the Dionysos client application.
// @securityDefinitions.basic BasicAuth

// @contact.name   API Support
// @contact.url    https://github.com/Brawdunoir/dionysos-server/issues

// @license.name  GNU General Public License v3.0
// @license.url   https://www.gnu.org/licenses/gpl-3.0.html
func main() {
	utils.InitAPI()

	// Gin initialization.
	router := routes.SetupRouter(gin.New())

	// Set VERSION in environment.
	os.Setenv("VERSION", VERSION)
	// Set dynamic documentation information.
	docs.SwaggerInfo.Version = VERSION
	docs.SwaggerInfo.BasePath = variables.BasePath

	// Start the server.
	l.Logger.Infof("Starting server on port %s", variables.Port)
	err := router.Run(":" + variables.Port)
	if err != nil {
		panic(err)
	}
}
