//nolint:typecheck
package routes

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the router
func SetupRouter(router *gin.Engine) *gin.Engine {
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
	}))

	userRouter := router.Group("/users")
	{
		userRouter.POST("", CreateUser)
		userRouter.GET("/:id", GetUser)
		userRouter.PATCH("/:id", UpdateUser)
		userRouter.DELETE("/:id", DeleteUser)
	}

	roomRouter := router.Group("/rooms")
	{
		roomRouter.POST("", CreateRoom)
		roomRouter.GET("/:id", GetRoom)
		roomRouter.PATCH("/:id", UpdateRoom)
		roomRouter.DELETE("/:id", DeleteRoom)
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
