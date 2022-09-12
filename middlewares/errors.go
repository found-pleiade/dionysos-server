package middlewares

import (
	"strings"

	utils "github.com/Brawdunoir/dionysos-server/utils/routes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {

			for _, ginErr := range c.Errors {
				logger.Error(ginErr.Error(), ginErr.Meta)
			}

			// status -1 doesn't overwrite existing status code
			lastErr := c.Errors[len(c.Errors)-1].Error()
			c.JSON(-1, utils.CreateErrorResponse(strings.ToUpper(lastErr[0:1])+lastErr[1:]))
		}
	}
}
