//nolint:typecheck
package routes

import (
	"time"

	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/Brawdunoir/dionysos-server/middlewares"
	"github.com/Brawdunoir/dionysos-server/utils"
	l "github.com/Brawdunoir/dionysos-server/utils/logger"
	"github.com/Brawdunoir/dionysos-server/variables"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var cacheStore persist.CacheStore

// Database pointer that will be used in the routes.
var db *gorm.DB

// SetupRouter sets up the router
func SetupRouter(router *gin.Engine) *gin.Engine {
	// Get the database connection.
	db = database.GetDB()
	if db == nil {
		l.Logger.Fatal("Failed to connect to database")
	}

	// Connect to Redis client or create a local redis.
	if variables.RedisHost != "" {
		redisURL, err := redis.ParseURL(variables.RedisHost)
		if err != nil {
			l.Logger.Error("Cannot connect to redis", err)
		} else {
			cacheStore = persist.NewRedisStore(redis.NewClient(redisURL))
		}
	} else {
		cacheStore = persist.NewMemoryStore(5 * time.Minute)
	}

	// Setup the routes.
	r := router.Group(variables.BasePath)
	{
		// Global middlewares.
		r.Use(
			gin.LoggerWithWriter(gin.DefaultWriter, "/healthz"),
			gin.Recovery(),
			middlewares.Options(),
		)
		// Public routes.
		r.GET("/healthz", Healthz)
		r.GET("/version", GetVersion)
		r.GET("/doc/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		r.POST("/users", CreateUser)

		// Add authentication middleware to the following routes.
		r.Use(middlewares.Authentication(db, l.Logger))

		userRouter := r.Group("/users")
		{
			userRouter.GET("/:id", cache.CacheByRequestURI(cacheStore, 5*time.Minute), GetUser)

			userRouter.Use(middlewares.InvalidateCacheURI(cacheStore, l.Logger))

			userRouter.PATCH("/:id", UpdateUser)
			userRouter.DELETE("/:id", DeleteUser)
		}

		roomRouter := r.Group("/rooms")
		{
			roomRouter.POST("", CreateRoom)

			roomRouter.Use(middlewares.RetrieveRoom(l.Logger, db))

			roomRouter.GET("/:id/stream", utils.HeadersSSE, StreamRoom)
			roomRouter.GET("/:id", cache.CacheByRequestURI(cacheStore, 5*time.Minute), GetRoom)

			roomRouter.Use(middlewares.InvalidateCacheURI(cacheStore, l.Logger))

			roomRouter.PATCH("/:id", UpdateRoom)
			roomRouter.PATCH("/:id/connect", ConnectUserToRoom)
			roomRouter.PATCH("/:id/disconnect", DisconnectUserFromRoom)
			roomRouter.PATCH("/:id/kick/:userid", KickUserFromRoom)
		}
	}

	router.Use(middlewares.ErrorHandler(l.Logger))
	return router
}
