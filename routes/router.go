//nolint:typecheck
package routes

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/Brawdunoir/dionysos-server/utils"
	routes "github.com/Brawdunoir/dionysos-server/utils/routes"
	"github.com/Brawdunoir/dionysos-server/variables"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var redisStore *persist.RedisStore

// Database pointer that will be used in the routes.
var db *gorm.DB = database.DB

// SetupRouter sets up the router
func SetupRouter(router *gin.Engine) *gin.Engine {
	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/healthz"),
		gin.Recovery(),
		options(),
	)

	// Connect to Redis.
	redisURL, err := redis.ParseURL(variables.RedisHost)
	if err != nil {
		log.Println("Cannot connect to redis", err)
	} else {
		redisStore = persist.NewRedisStore(redis.NewClient(redisURL))
	}

	r := router.Group(variables.BasePath)
	{
		// We should not use the authentication middleware on these routes.
		r.GET("/healthz", Healthz)
		r.GET("/version", GetVersion)
		r.GET("/doc/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		r.POST("/users", CreateUser)

		r.Use(authentication)

		userRouter := r.Group("/users")
		{
			if redisStore != nil {
				userRouter.GET("/:id", cache.CacheByRequestURI(redisStore, 60*time.Minute), GetUser)
				userRouter.PATCH("/:id", invalidateCacheURI("users"), UpdateUser)
				userRouter.DELETE("/:id", invalidateCacheURI("users"), DeleteUser)
			} else {
				userRouter.GET("/:id", GetUser)
				userRouter.PATCH("/:id", UpdateUser)
				userRouter.DELETE("/:id", DeleteUser)
			}
		}

		roomRouter := r.Group("/rooms")
		{
			roomRouter.POST("", CreateRoom)
			roomRouter.Use(retrieveRoom)
			roomRouter.GET("/:id/stream", utils.HeadersSSE, StreamRoom)
			if redisStore != nil {
				roomRouter.GET("/:id", cache.CacheByRequestURI(redisStore, 60*time.Minute), GetRoom)
				roomRouter.PATCH("/:id", invalidateCacheURI("rooms"), UpdateRoom)
				roomRouter.PATCH("/:id/connect", invalidateCacheURI("rooms"), ConnectUserToRoom)
				roomRouter.PATCH("/:id/disconnect", invalidateCacheURI("rooms"), DisconnectUserFromRoom)
				roomRouter.PATCH("/:id/kick/:userid", invalidateCacheURI("rooms"), KickUserFromRoom)
			} else {
				roomRouter.GET("/:id", GetRoom)
				roomRouter.PATCH("/:id", UpdateRoom)
				roomRouter.PATCH("/:id/connect", ConnectUserToRoom)
				roomRouter.PATCH("/:id/disconnect", DisconnectUserFromRoom)
				roomRouter.PATCH("/:id/kick/:userid", KickUserFromRoom)
			}
		}
	}

	return router
}

// Middleware to authenticate users.
func authentication(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()
	// Extract the id and password from the request Authorization header.
	id, password, ok := c.Request.BasicAuth()
	if ok {
		// Get user from database.
		var user models.User
		err := db.WithContext(ctx).First(&user, id).Error

		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found during authentication"})
			log.Printf("Failed to find user: %v", err)
			return
		}
		// Calculate SHA-256 hashes for the provided and expected passwords.
		passwordHash := sha256.Sum256([]byte(password))
		expectedPasswordHash := sha256.Sum256([]byte(user.Password))

		// Use the subtle.ConstantTimeCompare() function to avoid leaking information.
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		if passwordMatch || variables.Environment == variables.ENVIRONMENT_DEVELOPMENT {
			c.Set(variables.USER_CONTEXT_KEY, user)
			c.Next()
			return
		}
	}

	// If the Authentication header is not present, is invalid, or thepassword is wrong, then
	// set a WWW-Authenticate header to inform the client that we expect them
	// to use basic authentication and send a 401 Unauthorized response.
	c.Header("WWW-Authenticate", `Basic id:password charset="UTF-8"`)
	c.AbortWithStatusJSON(http.StatusUnauthorized, routes.CreateErrorResponse("User not authorized"))
}

// Middleware for CORS requests.
func options() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Allow", "GET,POST,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "X-Requested-With, Authorization, Origin, Content-Type, Accept")
		c.Header("Content-Type", "application/json,text/event-stream")
		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusOK)
		}
	}
}

// Invalidate cache for a given URI.
func invalidateCacheURI(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// Get the response code.
		code := c.Writer.Status()
		if code >= 200 && code < 300 && redisStore != nil {
			err := redisStore.Delete(resource + "/" + c.Param("id"))
			if err != nil {
				log.Printf("Failed to invalidate cache: %v", err)
			}
		}
	}
}

func retrieveRoom(c *gin.Context) {
	var room models.Room

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Failed to convert room ID: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.CreateErrorResponse("Invalid room ID"))
		return
	}

	err = room.GetRoom(c, db, id)
	if err != nil {
		log.Printf("Failed to find document: %v", err)
		c.AbortWithStatusJSON(http.StatusNotFound, routes.CreateErrorResponse("Room not found"))
		return
	}
	c.Set(variables.ROOM_CONTEXT_KEY, room)
	c.Next()
}
