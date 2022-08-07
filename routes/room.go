//nolint:typecheck
package routes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Brawdunoir/dionysos-server/models"
	utils "github.com/Brawdunoir/dionysos-server/utils"
	utilsRoutes "github.com/Brawdunoir/dionysos-server/utils/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateRoom creates a room in the database
func CreateRoom(c *gin.Context) {
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	user, err := utilsRoutes.ExtractUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in context"})
		log.Printf("Failed to extract user from context: %v", err)
		return
	}

	room.OwnerID = user.ID
	room.UsersID = append(room.UsersID, user.ID)

	err = db.WithContext(ctx).Create(&room).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not created"})
		log.Printf("Failed to create document: %v", err)
		return
	}

	c.JSON(http.StatusCreated, utilsRoutes.CreateResponse{URI: "/rooms/" + fmt.Sprint(room.ID)})
}

// GetRoom returns a room from the database
func GetRoom(c *gin.Context) {
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		log.Printf("Failed to convert room ID: %v", err)
	}

	err = db.WithContext(ctx).First(&room, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		log.Printf("Failed to find document: %v", err)
		return
	}

	c.JSON(http.StatusOK, room)
}

// UpdateRoom updates a room in the database
func UpdateRoom(c *gin.Context) {
	var roomUpdate models.RoomUpdate
	var patchedRoom models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		log.Printf("Failed to convert room ID: %v", err)
		return
	}

	// Test if data is valid
	if err := c.ShouldBindJSON(&roomUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	err = db.WithContext(ctx).First(&patchedRoom, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			log.Printf("Failed to find document: %v", err)
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not updated"})
			log.Printf("Failed to modify document: %v", err)
			return
		}
	} else {
		// Check if requester is the owner of the room
		err := utilsRoutes.AssertUser(c, int(patchedRoom.OwnerID))
		if err != nil {
			log.Printf("Error when asserting user: %v", err)
			return
		}

		err = db.WithContext(ctx).Model(&patchedRoom).Updates(roomUpdate.ToRoom()).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not modified"})
			log.Printf("Failed to modify document: %v", err)
			return
		}
	}

	c.JSON(http.StatusNoContent, nil)
}

// ConnectUserToRoom connects a user to a room in the database
func ConnectUserToRoom(c *gin.Context) {
	var user *models.User
	var room models.Room
	var roomUpdate models.RoomUpdate

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	user, err := utilsRoutes.ExtractUserFromContext(c)
	if err != nil {
		log.Printf("Failed to extract user from context: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user from context"})
		return
	}

	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		log.Printf("Failed to convert room ID: %v", err)
		return
	}

	err = db.First(&room, roomID).Error
	if err != nil {
		log.Printf("Failed to find document: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to find document"})
		return
	}

	roomUpdate.UsersID = append(roomUpdate.UsersID, user.ID)
	err = db.WithContext(ctx).Model(&room).Updates(roomUpdate.ToRoom()).Error
	if err != nil {
		log.Printf("Failed to modify document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to modify document"})
		return
	}

	c.JSON(http.StatusOK, room)
}

// DisconnectUserFromRoom disconnects a user from a room in the database
func DisconnectUserFromRoom(c *gin.Context) {
	var user *models.User
	var room models.Room
	var roomUpdate models.RoomUpdate

	_, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	user, err := utilsRoutes.ExtractUserFromContext(c)
	if err != nil {
		log.Printf("Failed to extract user from context: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user from context"})
		return
	}

	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		log.Printf("Failed to convert room ID: %v", err)
		return
	}

	err = db.First(&room, roomID).Error
	if err != nil {
		log.Printf("Failed to find document: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to find document"})
		return
	}

	// Remove user from the connected users list of the room
	roomUpdate.UsersID, err = utils.RemoveUintFromSlice(roomUpdate.UsersID, user.ID)
	if err != nil {
		log.Printf("User not connected to room: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from room"})
		return
	}

	// Check if user is owner to know if we give ownership to a new user
	if room.OwnerID == user.ID {
		if len(roomUpdate.UsersID) > 0 {
			err = utilsRoutes.ChangeOwner(roomUpdate, db)
			if err != nil {
				log.Printf("Failed to change owner: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change owner"})
				return
			}

		} else {
			err = utilsRoutes.DeleteRoom(room.ID, db)
			if err == nil {
				log.Printf("Room has been deleted")
			} else {
				log.Printf("Failed to delete room: %v", err)
			}
		}
	}

	// Delete user
	c.Set("id", user.ID)
	DeleteUser(c)

	c.JSON(http.StatusOK, room)
}
