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

	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/Brawdunoir/dionysos-server/variables"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the router
func SetupRouter(router *gin.Engine) *gin.Engine {
	basePath := variables.BasePath

	router.Use(options)

	r := router.Group(basePath)
	{
		// We should not use the basicAuth middleware for the /users endpoint because it is used to create the user.
		r.POST("/users", CreateUser)

		userRouter := r.Group("/users", basicAuth)
		{
			userRouter.GET("/:id", GetUser)
			userRouter.PATCH("/:id", UpdateUser)
			userRouter.DELETE("/:id", DeleteUser)
		}

		roomRouter := r.Group("/rooms", basicAuth)
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

// Middleware to authenticate users.
func basicAuth(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()
	// Extract the id and password from the request Authorization header.
	id, password, ok := c.Request.BasicAuth()
	if ok {
		// Get user from database.
		var user models.User
		err := db.WithContext(ctx).First(&user, id).Error

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			log.Printf("Failed to find document: %v", err)
			return
		}
		// Calculate SHA-256 hashes for the provided and expected passwords.
		passwordHash := sha256.Sum256([]byte(password))
		expectedPasswordHash := sha256.Sum256([]byte(user.Password))

		// Use the subtle.ConstantTimeCompare() function to avoid leaking information.
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		if passwordMatch {
			c.Set("id", id)
			c.Next()
			return
		}
	}

	// If the Authentication header is not present, is invalid, or the
	// password is wrong, then set a WWW-Authenticate
	// header to inform the client that we expect them to use basic
	// authentication and send a 401 Unauthorized response.
	c.Header("WWW-Authenticate", `Basic id:password charset="UTF-8"`)
	c.AbortWithStatus(http.StatusUnauthorized)
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
