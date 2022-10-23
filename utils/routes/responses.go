package utils

import (
	"errors"
	"net/http"

	e "github.com/Brawdunoir/dionysos-server/utils/errors"
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

// AssertUser compares the ID of the authenticated user in context and the ID of the room owner.
// It returns an error if the user is not the owner of the room.
// It also sets the JSON response so caller only needs to return if an error is returned.
func AssertUserIsOwner(c *gin.Context, expectedID uint64) error {
	user, err := ExtractUserFromContext(c)
	if err != nil {
		c.Error(err).SetMeta("AssertUser.ExtractUserFromContext")
		c.AbortWithError(http.StatusInternalServerError, e.UserNotInContext{}).SetMeta("AssertUser.ExtractUserFromContext")
		return err
	}

	if user.ID != expectedID {
		err := errors.New("user is not the owner of the room, not authorized")
		c.Error(err).SetMeta("AssertUser.CompareID")
		c.AbortWithError(http.StatusUnauthorized, e.UserNotAuthorized{}).SetMeta("AssertUser.CompareID")
		return err
	}

	return nil
}

func CreateErrorResponse(error string) *ErrorResponse {
	return &ErrorResponse{
		Error: error,
	}
}
