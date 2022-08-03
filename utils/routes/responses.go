package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateResponse struct {
	URI      string `json:"uri"`
	Password string `json:"password,omitempty"`
}

func AssertUser(c *gin.Context, expectedID int) error {
	user, err := ExtractUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context. Has it been set in the middleware?"})
		return err
	}

	if user.ID != uint(expectedID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return errors.New("user not authorized")
	}

	return nil
}
