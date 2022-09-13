package middlewares

import (
	"net/http"
	"strconv"

	"github.com/Brawdunoir/dionysos-server/models"
	e "github.com/Brawdunoir/dionysos-server/utils/errors"
	"github.com/Brawdunoir/dionysos-server/variables"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RetrieveRoom retrieves the room from the database and adds it to the context.
func RetrieveRoom(logger *zap.SugaredLogger, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var room models.Room

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.Error(err).SetMeta("RetrieveRoom.ParseUint")
			c.AbortWithError(http.StatusBadRequest, e.InvalidID{}).SetMeta("RetrieveRoom.ParseUint")
			return
		}

		err = room.GetRoom(c, db, id)
		if err != nil {
			c.Error(err).SetMeta("RetrieveRoom.GetRoom")
			c.AbortWithError(http.StatusNotFound, e.RoomNotFound{}).SetMeta("RetrieveRoom.GetRoom")
			return
		}
		c.Set(variables.ROOM_CONTEXT_KEY, room)
		c.Next()
	}
}
