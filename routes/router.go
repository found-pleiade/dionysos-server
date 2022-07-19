//nolint:typecheck
package routes

import (
	"net/http"
	"os"

	"github.com/Brawdunoir/dionysos-server/constants"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the router
func SetupRouter(router *gin.Engine) *gin.Engine {
	basePath := constants.BasePath

	router.Use(options)

	r := router.Group(basePath)
	{
		userRouter := r.Group("/users")
		{
			userRouter.POST("", CreateUser)
			userRouter.GET("/:id", GetUser)
			userRouter.PATCH("/:id", UpdateUser)
			userRouter.DELETE("/:id", DeleteUser)
		}

		roomRouter := r.Group("/rooms")
		{
			roomRouter.POST("", CreateRoom)
			roomRouter.GET("/:id", GetRoom)
			roomRouter.PATCH("/:id", UpdateRoom)
			roomRouter.DELETE("/:id", DeleteRoom)
		}
		r.GET("/version", func(c *gin.Context) {
			var version string
			if version = os.Getenv("VERSION"); version == "" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Version not set"})
			}
			c.String(http.StatusOK, version)
		})
	}
	return router
}

// Middleware for CORS requests.
func options(c *gin.Context) {
	c.Header("Allow", "GET,POST,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "X-Requested-With, Authorization, Origin, Content-Type, Accept")
	c.Header("Content-Type", "application/json")
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
