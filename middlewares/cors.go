package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware for CORS requests.
// It intercepts OPTIONS requests and adds the necessary headers.
func Options() gin.HandlerFunc {
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
