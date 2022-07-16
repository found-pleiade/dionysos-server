//nolint:typecheck
package routes

import (
	"github.com/Brawdunoir/dionysos-server/constants"
	docs "github.com/Brawdunoir/dionysos-server/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter sets up the router
func SetupRouter(router *gin.Engine) *gin.Engine {

	basePath := constants.BasePath

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
	}))

	v := router.Group(basePath)
	{
		// Users
		userRouter := v.Group("/users")
		{
			userRouter.POST("/", CreateUser)
			userRouter.GET("/:id", GetUser)
			userRouter.PATCH("/:id", UpdateUser)
			userRouter.DELETE("/:id", DeleteUser)
		}

		// Rooms
		roomRouter := v.Group("/rooms")
		{
			roomRouter.POST("/", CreateRoom)
			roomRouter.GET("/:id", GetRoom)
			roomRouter.PATCH("/:id", UpdateRoom)
			roomRouter.DELETE("/:id", DeleteRoom)
		}
		// Version
		v.GET("/version", GetVersion)
	}

	// Docs
	docs.SwaggerInfo.BasePath = basePath
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
