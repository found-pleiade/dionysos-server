package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, ginErr := range c.Errors {
			logger.Error(ginErr.Error(), ginErr.Meta)
		}

		// status -1 doesn't overwrite existing status code
		c.JSON(-1, gin.H{"error": "Internal server error"})

	}
}
