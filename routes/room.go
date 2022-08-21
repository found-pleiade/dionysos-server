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
	routes "github.com/Brawdunoir/dionysos-server/utils/routes"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

// CreateRoom godoc
// @Summary      Creates a room.
// @Tags         Rooms
// @Security 		 BasicAuth
// @Accept       json
// @Produce      json
// @Param 			 room	body	models.RoomUpdate	true	"Room object"
// @Success      201 {object} utils.CreateResponse "Room created"
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "Invalid user in auth method"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms [post]
func CreateRoom(c *gin.Context) {
	var r models.RoomUpdate

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse(err.Error()))
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	user, err := routes.ExtractUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusNotFound, routes.CreateErrorResponse("User not found in context"))
		log.Printf("Failed to extract user from context: %v", err)
		return
	}

	room := r.ToRoom()

	room.ID, err = utils.UUIDGenerator.NextID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("Room not created"))
		log.Printf("Failed to create document: %v", err)
		return
	}

	room.OwnerID = user.ID
	room.Users = append(room.Users, user)

	err = db.WithContext(ctx).Create(&room).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("Room not created"))
		log.Printf("Failed to create document: %v", err)
		return
	}

	c.JSON(http.StatusCreated, routes.CreateResponse{URI: "/rooms/" + fmt.Sprint(room.ID)})
}

// GetRoom godoc
// @Summary      Gets a room.
// @Tags         Rooms
// @Security 		 BasicAuth
// @Produce      json
// @Param 			 id	path	int	true	"Room ID"
// @Success      200 {object} models.Room
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "Room not found or invalid user in auth method"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms/{id} [get]
func GetRoom(c *gin.Context) {
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse("Invalid room ID"))
		log.Printf("Failed to convert room ID: %v", err)
	}

	err = room.GetRoom(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, routes.CreateErrorResponse("Room not found"))
		log.Printf("Failed to find document: %v", err)
		return
	}

	c.JSON(http.StatusOK, room)
}

// UpdateRoom godoc
// @Summary      Updates a room.
// @Tags         Rooms
// @Security 		 BasicAuth
// @Accept       json
// @Produce      json
// @Param 			 id		path	int								true	"Room ID"
// @Param 			 room body 	models.RoomUpdate true	"Room object"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "Room not found or invalid user in auth method"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms/{id} [patch]
func UpdateRoom(c *gin.Context) {
	var r models.RoomUpdate
	var patchedRoom models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse("Invalid room ID"))
		log.Printf("Failed to convert room ID: %v", err)
		return
	}

	// Test if data is valid
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse(err.Error()))
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	err = db.WithContext(ctx).First(&patchedRoom, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, routes.CreateErrorResponse("Room not found"))
			return
		} else {
			c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("Room not updated"))
			log.Printf("Failed to modify document: %v", err)
			return
		}
	} else {
		// Check if requester is the owner of the room
		err := routes.AssertUser(c, patchedRoom.OwnerID)
		if err != nil {
			log.Printf("Error when asserting user: %v", err)
			return
		}

		err = db.WithContext(ctx).Model(&patchedRoom).Updates(r.ToRoom()).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("Room not modified"))
			log.Printf("Failed to modify document: %v", err)
			return
		}
	}

	c.JSON(http.StatusNoContent, nil)
}

// ConnectUserToRoom godoc
// @Summary      Connects a user to a room.
// @Tags         Rooms
// @Security 		 BasicAuth
// @Produce      json
// @Param 			 id	path	int	true	"Room ID"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "Room not found or invalid user in auth method"
// @Failure      409 {object} utils.ErrorResponse "User already in room"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms/{id}/connect [patch]
func ConnectUserToRoom(c *gin.Context) {
	var user models.User
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	user, err := routes.ExtractUserFromContext(c)
	if err != nil {
		log.Printf("Failed to extract user from context: %v", err)
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("Failed to extract user from context"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse("Invalid room ID"))
		log.Printf("Failed to convert room ID: %v", err)
		return
	}

	err = room.GetRoom(ctx, db, id)
	if err != nil {
		log.Printf("Failed to find document: %v", err)
		c.JSON(http.StatusNotFound, routes.CreateErrorResponse("Room not found"))
		return
	}

	// Assert user is not already in the room.
	if slices.Contains(room.Users, user) {
		log.Printf("User already in room")
		c.JSON(http.StatusConflict, routes.CreateErrorResponse("User already in room"))
		return
	}

	room.Users = append(room.Users, user)

	err = db.WithContext(ctx).Save(&room).Error
	if err != nil {
		log.Printf("Failed to modify document: %v", err)
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("Room not modified"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// DisconnectUserFromRoom godoc
// @Summary      Disconnects a user from a room.
// @Tags         Rooms
// @Security 		 BasicAuth
// @Produce      json
// @Param 			 id	path	int	true	"Room ID"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "Room not found or invalid user in auth method"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms/{id}/disconnect [patch]
func DisconnectUserFromRoom(c *gin.Context) {
	var user models.User
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	user, err := routes.ExtractUserFromContext(c)
	if err != nil {
		log.Printf("Failed to extract user from context: %v", err)
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("Failed to extract user from context"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Failed to convert room ID: %v", err)
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse("Invalid room ID"))
		return
	}

	err = room.GetRoom(ctx, db, id)
	if err != nil {
		log.Printf("Failed to find document: %v", err)
		c.JSON(http.StatusNotFound, routes.CreateErrorResponse("Room not found"))
		return
	}

	if !slices.Contains(room.Users, user) {
		log.Printf("User not connected to room: %v", err)
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse("User not in room"))
		return
	}

	// Remove user from the connected users list of the room
	err = room.RemoveUser(ctx, db, &user)
	if err != nil {
		log.Printf("Failed to remove user from room: %v", err)
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse("Failed to remove user from room"))
		return
	}

	// We want to delete an empty room and keep an owner at every instant.
	if len(room.Users) == 0 {
		result := db.WithContext(ctx).Delete(&room)
		if result.Error != nil {
			log.Printf("Failed to delete document: %v", err)
			c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("Room not deleted"))
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
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("Room not modified"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
