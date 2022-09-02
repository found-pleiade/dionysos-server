//nolint:typecheck
package routes

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/Brawdunoir/dionysos-server/utils"
	routes "github.com/Brawdunoir/dionysos-server/utils/routes"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

// CreateRoom godoc
// @Summary      Creates a room.
// @Tags         Rooms
// @Security     BasicAuth
// @Accept       json
// @Produce      json
// @Param        room body models.RoomUpdate true "Room object"
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
// @Security     BasicAuth
// @Produce      json
// @Param        id path int true "Room ID"
// @Success      200 {object} models.Room
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "Room not found or invalid user in auth method"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms/{id} [get]
func GetRoom(c *gin.Context) {
	var room models.Room

	_, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	room, err := routes.ExtractRoomFromContext(c)
	if err != nil {
		log.Printf("Failed to extract room from context: %v", err)
		c.JSON(http.StatusNotFound, routes.CreateErrorResponse("Room not found"))
		return
	}

	c.JSON(http.StatusOK, room)
}

// UpdateRoom godoc
// @Summary      Updates a room.
// @Tags         Rooms
// @Security     BasicAuth
// @Accept       json
// @Produce      json
// @Param        id   path int               true "Room ID"
// @Param        room body models.RoomUpdate true "Room object"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "Room not found or invalid user in auth method"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms/{id} [patch]
func UpdateRoom(c *gin.Context) {
	var r models.RoomUpdate

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	room, err := routes.ExtractRoomFromContext(c)
	if err != nil {
		log.Printf("Failed to extract room from context: %v", err)
		c.JSON(http.StatusNotFound, routes.CreateErrorResponse("Room not found"))
		return
	}

	// Test if data is valid
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse(err.Error()))
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	// Check if requester is the owner of the room
	err = routes.AssertUser(c, room.OwnerID)
	if err != nil {
		log.Printf("Error when asserting user: %v", err)
		return

	}

	err = db.WithContext(ctx).Model(&room).Updates(r.ToRoom()).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("Room not modified"))
		log.Printf("Failed to modify document: %v", err)
		return
	}

	stream, ok := listStreamsRoom[room.ID]
	if ok {
		stream.distribute(Message{Event: "roomNameUpdated", Data: room.Name})
	} else {
		log.Printf("Problem with stream of room %v", room.ID)
	}

	c.JSON(http.StatusNoContent, nil)
}

// ConnectUserToRoom godoc
// @Summary      Connects a user to a room.
// @Tags         Rooms
// @Security     BasicAuth
// @Produce      json
// @Param        id path int true "Room ID"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "Room not found or invalid user in auth method"
// @Failure      409 {object} utils.ErrorResponse "User already in room"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms/{id}/connect [patch]
func ConnectUserToRoom(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	user, err := routes.ExtractUserFromContext(c)
	if err != nil {
		log.Printf("Failed to extract user from context: %v", err)
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("Failed to extract user from context"))
		return
	}

	room, err := routes.ExtractRoomFromContext(c)
	if err != nil {
		log.Printf("Failed to extract room from context: %v", err)
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

	stream, ok := listStreamsRoom[room.ID]
	if ok {
		stream.distribute(Message{Event: "userConnected", Data: user})
	} else {
		log.Printf("Problem with stream of room %v", room.ID)
	}

	c.JSON(http.StatusNoContent, nil)
}

// DisconnectUserFromRoom godoc
// @Summary      Disconnects a user from a room.
// @Tags         Rooms
// @Security     BasicAuth
// @Produce      json
// @Param        id path int true "Room ID"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "Room not found or invalid user in auth method"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms/{id}/disconnect [patch]
func DisconnectUserFromRoom(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	user, err := routes.ExtractUserFromContext(c)
	if err != nil {
		log.Printf("Failed to extract user from context: %v", err)
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("Failed to extract user from context"))
		return
	}

	room, err := routes.ExtractRoomFromContext(c)
	if err != nil {
		log.Printf("Failed to extract room from context: %v", err)
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

	stream, ok := listStreamsRoom[room.ID]
	if ok {
		stream.distribute(Message{Event: "userDisconnected", Data: user.ID})
	} else {
		log.Printf("Problem with stream of room %v", room.ID)
	}

	c.JSON(http.StatusNoContent, nil)
}

// StreamRoom WIP.
func StreamRoom(c *gin.Context) {
	var stream Stream
	room, err := routes.ExtractRoomFromContext(c)
	if err != nil {
		log.Printf("Failed to extract room from context: %v", err)
		c.AbortWithStatusJSON(http.StatusNotFound, routes.CreateErrorResponse("Room not found"))
	}

	user, err := routes.ExtractUserFromContext(c)
	if err != nil {
		log.Printf("Failed to extract user from context: %v", err)
		c.JSON(http.StatusNotFound, routes.CreateErrorResponse("User not found"))
		return
	}

	stream, ok := listStreamsRoom[room.ID]
	if !ok {
		stream.create()
		listStreamsRoom[room.ID] = stream
	}

	stream.addSub(user.ID)

	close := c.Stream(func(w io.Writer) bool {
		if msg, ok := <-stream.ClientChan[user.ID]; ok {
			c.SSEvent(msg.Event, msg.Data)
			return true
		}
		return false
	})
	if close {
		stream.delSub(user.ID)
	}
}
