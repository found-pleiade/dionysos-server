package middlewares

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"time"

	"github.com/Brawdunoir/dionysos-server/models"
	e "github.com/Brawdunoir/dionysos-server/utils/errors"
	"github.com/Brawdunoir/dionysos-server/variables"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Middleware to authenticate users.
// It also places the user in the context for later use.
func Authentication(db *gorm.DB, logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
		defer cancelCtx()
		// Extract the id and password from the request Authorization header.
		id, password, ok := c.Request.BasicAuth()
		if ok {
			// Get user from database.
			var user models.User
			err := db.WithContext(ctx).First(&user, id).Error

			if err != nil {
				c.Error(err).SetMeta("Authentication.First")
				c.AbortWithError(http.StatusNotFound, e.UserNotFound{}).SetMeta("Authentication.First")
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
		c.AbortWithError(http.StatusUnauthorized, e.UserNotAuthorized{}).SetMeta("Authentication")
	}
}
