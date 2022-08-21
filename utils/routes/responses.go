package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateResponse struct {
	URI string `json:"uri"`

	// Only set for user creation
	Password string `json:"password,omitempty"`
}

func AssertUser(c *gin.Context, expectedID uint64) error {
	user, err := ExtractUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CreateErrorResponse("User not found in context. Has it been set in the middleware?"))
		return err
	}

	if user.ID != expectedID {
		c.JSON(http.StatusUnauthorized, CreateErrorResponse("User not authorized"))
		return errors.New("user not authorized")
	}

	return nil
}

func CreateErrorResponse(error string) *ErrorResponse {
	return &ErrorResponse{
		Error: error,
	}
}
