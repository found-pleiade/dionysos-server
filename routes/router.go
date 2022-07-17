//nolint:typecheck
package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the router
func SetupRouter(router *gin.Engine) *gin.Engine {
	router.Use(options)

	userRouter := router.Group("/users")
	{
		userRouter.POST("/", CreateUser)
		userRouter.GET("/:id/", GetUser)
		userRouter.PATCH("/:id/", UpdateUser)
		userRouter.DELETE("/:id/", DeleteUser)
	}

	roomRouter := router.Group("/rooms")
	{
		roomRouter.POST("/", CreateRoom)
		roomRouter.GET("/:id/", GetRoom)
		roomRouter.PATCH("/:id/", UpdateRoom)
		roomRouter.DELETE("/:id/", DeleteRoom)
	}

	router.GET("/version", func(c *gin.Context) {
		var version string
		if version = os.Getenv("VERSION"); version == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Version not set"})
		}
		c.String(http.StatusOK, version)
	})

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
