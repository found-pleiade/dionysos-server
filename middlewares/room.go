package middlewares

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Brawdunoir/dionysos-server/models"
	utils "github.com/Brawdunoir/dionysos-server/utils/routes"
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
			log.Printf("Failed to convert room ID: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid room ID"))
			return
		}

		err = room.GetRoom(c, db, id)
		if err != nil {
			log.Printf("Failed to find document: %v", err)
			c.AbortWithStatusJSON(http.StatusNotFound, utils.CreateErrorResponse("Room not found"))
			return
		}
		c.Set(variables.ROOM_CONTEXT_KEY, room)
		c.Next()
	}
}
