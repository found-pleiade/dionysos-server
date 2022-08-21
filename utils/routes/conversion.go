package utils

import (
	"errors"

	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/Brawdunoir/dionysos-server/variables"
	"github.com/gin-gonic/gin"
)

func ExtractUserFromContext(c *gin.Context) (models.User, error) {
	val, ok := c.Get(variables.USER_CONTEXT_KEY)
	if !ok {
		return models.User{}, errors.New("user not found in context")
	}
	user, ok := val.(models.User)
	if !ok {
		return models.User{}, errors.New("value found in context but it was not a user")
	}

	return user, nil
}

func ExtractRoomFromContext(c *gin.Context) (models.Room, error) {
	val, ok := c.Get(variables.ROOM_CONTEXT_KEY)
	if !ok {
		return models.Room{}, errors.New("room not found in context")
	}
	room, ok := val.(models.Room)
	if !ok {
		return models.Room{}, errors.New("value found in context but it was not a room")
	}

	return room, nil
}
