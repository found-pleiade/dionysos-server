package routes

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the router
func SetupRouter() *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))

	userRouter := router.Group("/users")
	{
		userRouter.POST("/", CreateUser)
		userRouter.GET("/:id", GetUser)
		userRouter.PATCH("/:id", UpdateUser)
		userRouter.DELETE("/:id", DeleteUser)
	}

	roomRouter := router.Group("/rooms")
	{
		roomRouter.POST("/", CreateRoom)
		roomRouter.GET("/:id", GetRoom)
		roomRouter.PATCH("/:id", UpdateRoom)
		roomRouter.DELETE("/:id", DeleteRoom)
	}

	router.GET("/version", func(ctx *gin.Context) {
		var version string
		if version = os.Getenv("VERSION"); version == "" {
			version = "unknown"
		}
		ctx.String(http.StatusOK, version)
	})

	return router
}
