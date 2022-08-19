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
	"github.com/Brawdunoir/dionysos-server/utils"
	utilsRoutes "github.com/Brawdunoir/dionysos-server/utils/routes"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

// CreateRoom creates a room in the database.
func CreateRoom(c *gin.Context) {
	var r models.RoomUpdate

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&r); err != nil {
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

	room := r.ToRoom()

	room.ID, err = utils.UUIDGenerator.NextID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not created"})
		log.Printf("Failed to create document: %v", err)
		return
	}

	room.OwnerID = user.ID
	room.Users = append(room.Users, user)

	err = db.WithContext(ctx).Create(&room).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not created"})
		log.Printf("Failed to create document: %v", err)
		return
	}

	c.JSON(http.StatusCreated, utilsRoutes.CreateResponse{URI: "/rooms/" + fmt.Sprint(room.ID)})
}

// GetRoom returns a room from the database.
func GetRoom(c *gin.Context) {
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		log.Printf("Failed to convert room ID: %v", err)
	}

	err = room.GetRoom(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		log.Printf("Failed to find document: %v", err)
		return
	}

	c.JSON(http.StatusOK, room)
}

// UpdateRoom updates a room in the database.
func UpdateRoom(c *gin.Context) {
	var r models.RoomUpdate
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
	if err := c.ShouldBindJSON(&r); err != nil {
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
		err := utilsRoutes.AssertUser(c, patchedRoom.OwnerID)
		if err != nil {
			log.Printf("Error when asserting user: %v", err)
			return
		}

		err = db.WithContext(ctx).Model(&patchedRoom).Updates(r.ToRoom()).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not modified"})
			log.Printf("Failed to modify document: %v", err)
			return
		}
	}

	c.JSON(http.StatusNoContent, nil)
}

// ConnectUserToRoom connects a user to a room in the database.
func ConnectUserToRoom(c *gin.Context) {
	var user models.User
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	user, err := utilsRoutes.ExtractUserFromContext(c)
	if err != nil {
		log.Printf("Failed to extract user from context: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user from context"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		log.Printf("Failed to convert room ID: %v", err)
		return
	}

	err = room.GetRoom(ctx, db, id)
	if err != nil {
		log.Printf("Failed to find document: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// Assert user is not already in the room.
	if slices.Contains(room.Users, user) {
		log.Printf("User already in room")
		c.JSON(http.StatusConflict, gin.H{"error": "User already in room"})
		return
	}

	room.Users = append(room.Users, user)

	err = db.WithContext(ctx).Save(&room).Error
	if err != nil {
		log.Printf("Failed to modify document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not modified"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// DisconnectUserFromRoom disconnects a user from a room in the database.
func DisconnectUserFromRoom(c *gin.Context) {
	var user models.User
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	user, err := utilsRoutes.ExtractUserFromContext(c)
	if err != nil {
		log.Printf("Failed to extract user from context: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user from context"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Failed to convert room ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	err = room.GetRoom(ctx, db, id)
	if err != nil {
		log.Printf("Failed to find document: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	if !slices.Contains(room.Users, user) {
		log.Printf("User not connected to room: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not in room"})
		return
	}

	// Remove user from the connected users list of the room
	err = room.RemoveUser(ctx, db, &user)
	if err != nil {
		log.Printf("Failed to remove user from room: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to remove user from room"})
		return
	}

	// We want to delete an empty room and keep an owner at every instant.
	if len(room.Users) == 0 {
		result := db.WithContext(ctx).Delete(&room)
		if result.Error != nil {
			log.Printf("Failed to delete document: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not deleted"})
			return
		}
		log.Printf("Room %v deleted", room.ID)
		c.JSON(http.StatusNoContent, nil)
		return
	} else if room.OwnerID == user.ID {
		room.OwnerID = room.Users[0].ID
	}

	err = db.WithContext(ctx).Save(&room).Error
	if err != nil {
		log.Printf("Failed to modify document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not modified"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
