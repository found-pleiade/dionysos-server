//nolint:typecheck
package routes

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/Brawdunoir/dionysos-server/variables"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var redisStore *persist.RedisStore

// Database pointer that will be used in the routes.
var db *gorm.DB = database.DB

// SetupRouter sets up the router
func SetupRouter(router *gin.Engine) *gin.Engine {
	basePath := variables.BasePath

	router.Use(options)
	router.Use(gin.Recovery())

	// Connect to Redis.
	redisURL, err := redis.ParseURL(variables.RedisHost)
	if err != nil {
		log.Println("Cannot connect to redis", err)
	} else {
		redisStore = persist.NewRedisStore(redis.NewClient(redisURL))
	}

	r := router.Group(basePath)
	{
		// We should not use the authentication middleware for the /users endpoint because the password is generated during the user creation.
		r.POST("/users", CreateUser)

		userRouter := r.Group("/users", authentication)
		{
			if redisStore != nil {
				userRouter.GET("/:id", cache.CacheByRequestURI(redisStore, 60*time.Minute), GetUser)
				userRouter.PATCH("/:id", invalidateCacheURI, UpdateUser)
				userRouter.DELETE("/:id", invalidateCacheURI, DeleteUser)
			} else {
				userRouter.GET("/:id", GetUser)
				userRouter.PATCH("/:id", UpdateUser)
				userRouter.DELETE("/:id", DeleteUser)
			}
		}

		roomRouter := r.Group("/rooms", authentication)
		{
			roomRouter.POST("", CreateRoom)
			if redisStore != nil {
				roomRouter.GET("/:id", cache.CacheByRequestURI(redisStore, 60*time.Minute), GetRoom)
				roomRouter.PATCH("/:id", invalidateCacheURI, UpdateRoom)
				roomRouter.PATCH("/:id/connect", invalidateCacheURI, ConnectUserToRoom)
				roomRouter.PATCH("/:id/disconnect", invalidateCacheURI, DisconnectUserFromRoom)
			} else {
				roomRouter.GET("/:id", GetRoom)
				roomRouter.PATCH("/:id", UpdateRoom)
				roomRouter.PATCH("/:id/connect", ConnectUserToRoom)
				roomRouter.PATCH("/:id/disconnect", DisconnectUserFromRoom)
			}
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
			log.Printf("Failed to find document: %v", err)
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
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
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

func invalidateCacheURI(c *gin.Context) {
	c.Next()
	code := c.Request.Response.StatusCode
	if !c.IsAborted() && code >= 200 && code < 300 && redisStore != nil {
		err := redisStore.Delete(c.Request.RequestURI)
		if err != nil {
			log.Printf("Failed to invalidate cache: %v", err)
		}
	}
}
