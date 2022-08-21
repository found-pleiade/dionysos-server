package main

import (
	"os"

	"github.com/Brawdunoir/dionysos-server/docs"
	"github.com/Brawdunoir/dionysos-server/routes"
	"github.com/Brawdunoir/dionysos-server/variables"
	"github.com/gin-gonic/gin"
)

// VERSION as a constant variable
const VERSION = "0.2.0"

// @title           Dionysos
// @description     API instance for dionysos client application.
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

	router := routes.SetupRouter(gin.Default())

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
