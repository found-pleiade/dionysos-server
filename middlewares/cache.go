package middlewares

import (
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Invalidate cache for a given URI.
func InvalidateCacheURI(cacheStore persist.CacheStore, logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.RequestURI

		// Get the response code.
		code := c.Writer.Status()
		if code >= 200 && code < 300 && cacheStore != nil {
			err := cacheStore.Delete(key)
			if err != nil {
				logger.Errorf("Failed to invalidate cache with key '%s': %v", key, err)
			} else {
				logger.Debugln("Cache deleted for key", key)
			}
		}
	}
}
