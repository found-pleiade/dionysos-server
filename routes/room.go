//nolint:typecheck
package routes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/Brawdunoir/dionysos-server/utils"
	e "github.com/Brawdunoir/dionysos-server/utils/errors"
	l "github.com/Brawdunoir/dionysos-server/utils/logger"
	routes "github.com/Brawdunoir/dionysos-server/utils/routes"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

// Keep track of all SSE channels that are currently on service.
var roomStreamsList = make(map[uint64]*utils.Stream)
var SSEMessage = utils.Message{Event: "roomUpdate"}

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
		c.Error(err).SetMeta("CreateRoom.ShouldBindJSON")
		c.AbortWithError(http.StatusBadRequest, e.FailJSONBind{}).SetMeta("CreateRoom.ShouldBindJSON")
		return
	}

	user, err := routes.ExtractUserFromContext(c)
	if err != nil {
		c.Error(err).SetMeta("CreateRoom.ExtractUserFromContext")
		c.AbortWithError(http.StatusInternalServerError, e.UserNotInContext{}).SetMeta("CreateRoom.ExtractUserFromContext")
		return
	}

	room := r.ToRoom()

	room.ID, err = utils.UUIDGenerator.NextID()
	if err != nil {
		c.Error(err).SetMeta("CreateRoom.UUIDGenerator.NextID")
		c.AbortWithError(http.StatusInternalServerError, e.RoomNotCreated{}).SetMeta("CreateRoom.UUIDGenerator.NextID")
		return
	}

	room.OwnerID = user.ID
	room.Users = append(room.Users, user)

	err = db.WithContext(ctx).Create(&room).Error
	if err != nil {
		c.Error(err).SetMeta("CreateRoom.Create")
		c.AbortWithError(http.StatusInternalServerError, e.RoomNotCreated{}).SetMeta("CreateRoom.Create")
		return
	}

	// Create a new SSE channel for the room.
	_ = utils.CreateStream(room.ID, roomStreamsList)

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
		c.Error(err).SetMeta("GetRoom.ExtractRoomFromContext")
		c.AbortWithError(http.StatusNotFound, e.RoomNotInContext{}).SetMeta("GetRoom.ExtractRoomFromContext")
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
		c.Error(err).SetMeta("UpdateRoom.ExtractRoomFromContext")
		c.AbortWithError(http.StatusInternalServerError, e.RoomNotInContext{}).SetMeta("UpdateRoom.ExtractRoomFromContext")
		return
	}

	// Test if data is valid
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Error(err).SetMeta("UpdateRoom.ShouldBindJSON")
		c.AbortWithError(http.StatusBadRequest, e.FailJSONBind{}).SetMeta("UpdateRoom.ShouldBindJSON")
		return
	}

	// Check if requester is the owner of the room
	err = routes.AssertUserIsOwner(c, room.OwnerID)
	if err != nil {
		return
	}

	err = db.WithContext(ctx).Model(&room).Updates(r.ToRoom()).Error
	if err != nil {
		c.Error(err).SetMeta("UpdateRoom.Updates")
		c.AbortWithError(http.StatusInternalServerError, e.RoomNotModified{}).SetMeta("UpdateRoom.Updates")
		return
	}

	stream, err := utils.GetStream(room.ID, roomStreamsList)
	if err != nil {
		l.Logger.Warnf("Failed to get stream: %v", err)
	} else {
		stream.Distribute(SSEMessage)
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
		c.Error(err).SetMeta("ConnectUserToRoom.ExtractUserFromContext")
		c.AbortWithError(http.StatusInternalServerError, e.UserNotInContext{}).SetMeta("ConnectUserToRoom.ExtractUserFromContext")
		return
	}

	room, err := routes.ExtractRoomFromContext(c)
	if err != nil {
		c.Error(err).SetMeta("ConnectUserToRoom.ExtractRoomFromContext")
		c.AbortWithError(http.StatusInternalServerError, e.RoomNotInContext{}).SetMeta("ConnectUserToRoom.ExtractRoomFromContext")
		return
	}

	// Assert user is not already in the room.
	if slices.Contains(room.Users, user) {
		c.AbortWithError(http.StatusConflict, e.UserAlreadyInRoom{}).SetMeta("ConnectUserToRoom.Contains")
		return
	}

	room.Users = append(room.Users, user)

	err = db.WithContext(ctx).Save(&room).Error
	if err != nil {
		c.Error(err).SetMeta("ConnectUserToRoom.Save")
		c.AbortWithError(http.StatusInternalServerError, e.RoomNotModified{}).SetMeta("ConnectUserToRoom.Save")
		return
	}

	stream, err := utils.GetStream(room.ID, roomStreamsList)
	if err != nil {
		l.Logger.Warnf("Failed to get stream: %v", err)
	} else {
		stream.Distribute(SSEMessage)
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
		c.Error(err).SetMeta("DisconnectUserFromRoom.ExtractUserFromContext")
		c.AbortWithError(http.StatusInternalServerError, e.UserNotInContext{}).SetMeta("DisconnectUserFromRoom.ExtractUserFromContext")
		return
	}

	room, err := routes.ExtractRoomFromContext(c)
	if err != nil {
		c.Error(err).SetMeta("DisconnectUserFromRoom.ExtractRoomFromContext")
		c.AbortWithError(http.StatusInternalServerError, e.RoomNotInContext{}).SetMeta("DisconnectUserFromRoom.ExtractRoomFromContext")
		return
	}

	// Remove user from the connected users list of the room
	err = room.RemoveUser(ctx, db, &user)
	if err != nil {
		c.Error(err).SetMeta("DisconnectUserFromRoom.RemoveUser")
		c.AbortWithError(http.StatusBadRequest, e.RoomNotModified{}).SetMeta("DisconnectUserFromRoom.RemoveUser")
		return
	}

	// We want to delete an empty room and keep an owner at every instant.
	if len(room.Users) == 0 {
		result := db.WithContext(ctx).Delete(&room)
		if result.Error != nil {
			c.Error(result.Error).SetMeta("DisconnectUserFromRoom.Delete")
			c.AbortWithError(http.StatusInternalServerError, e.RoomNotModified{}).SetMeta("DisconnectUserFromRoom.Delete")
			return
		}
		l.Logger.Infof("Room %v deleted", room.ID)
		c.JSON(http.StatusNoContent, nil)
		return
	} else if room.OwnerID == user.ID {
		room.OwnerID = room.Users[0].ID
	}

	err = db.WithContext(ctx).Save(&room).Error
	if err != nil {
		c.Error(err).SetMeta("DisconnectUserFromRoom.Save")
		c.AbortWithError(http.StatusInternalServerError, e.RoomNotModified{}).SetMeta("DisconnectUserFromRoom.Save")
		return
	}

	stream, err := utils.GetStream(room.ID, roomStreamsList)
	if err != nil {
		l.Logger.Warnf("Failed to get stream: %v", err)
	} else {
		stream.Distribute(SSEMessage)
	}

	c.JSON(http.StatusNoContent, nil)
}

// StreamRoom godoc
// @Summary      SSE stream of a room for any updates.
// @Description  This endpoint is used to subscribe to a SSE stream for a given room.
// @Description	 The stream will send an event when a room is updated.
// @Description  A room is updated when a user connects or disconnects from it, or when we have a owner change, and so on.
// @Tags         Rooms,SSE
// @Security     BasicAuth
// @Param        id path int true "Room ID"
// @Produce      text/event-stream
// @Success      200 "Send \"RoomUpdate\" event each time room is updated. Send 200 when stream is closed"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "Room not found or invalid user in auth method"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms/{id}/stream [get]
func StreamRoom(c *gin.Context) {
	room, err := routes.ExtractRoomFromContext(c)
	if err != nil {
		c.Error(err).SetMeta("StreamRoom.ExtractRoomFromContext")
		c.AbortWithError(http.StatusInternalServerError, e.RoomNotInContext{}).SetMeta("StreamRoom.ExtractRoomFromContext")
		return
	}

	user, err := routes.ExtractUserFromContext(c)
	if err != nil {
		c.Error(err).SetMeta("StreamRoom.ExtractUserFromContext")
		c.AbortWithError(http.StatusInternalServerError, e.UserNotInContext{}).SetMeta("StreamRoom.ExtractUserFromContext")
		return
	}

	stream, err := utils.GetStream(room.ID, roomStreamsList)
	if err != nil {
		c.Error(err).SetMeta("StreamRoom.GetStream")
		c.AbortWithError(http.StatusInternalServerError, e.StreamNotCreated{}).SetMeta("StreamRoom.GetStream")
		return
	}

	_ = stream.AddSub(user.ID)
	defer func() {
		err := stream.DelSub(user.ID)
		if err != nil {
			l.Logger.Warnf("Failed to delete sub: %v", err)
		}
	}()

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-stream.ClientChan[user.ID]; ok {
			c.SSEvent(msg.Event, msg.Data)
			return true
		}
		return false
	})
}

// KickUserFromRoom godoc
// @Summary      Kicks a user from a room.
// @Tags         Rooms
// @Security     BasicAuth
// @Produce      json
// @Param        id 	path int true "Room ID"
// @Param        userid path int true "User ID"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "Room not found or invalid user in auth method"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms/{id}/kick/{userid} [patch]
func KickUserFromRoom(c *gin.Context) {
	var user models.User
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	room, err := routes.ExtractRoomFromContext(c)
	if err != nil {
		c.Error(err).SetMeta("KickUserFromRoom.ExtractRoomFromContext")
		c.AbortWithError(http.StatusInternalServerError, e.RoomNotInContext{}).SetMeta("KickUserFromRoom.ExtractRoomFromContext")
		return
	}

	userID, err := strconv.ParseUint(c.Param("userid"), 10, 64)
	if err != nil {
		c.Error(err).SetMeta("KickUserFromRoom.ParseUint")
		c.AbortWithError(http.StatusBadRequest, e.InvalidID{}).SetMeta("KickUserFromRoom.ParseUint")
		return
	}

	err = db.WithContext(ctx).First(&user, userID).Error
	if err != nil {
		c.Error(err).SetMeta("KickUserFromRoom.First")
		c.AbortWithError(http.StatusNotFound, e.UserNotFound{}).SetMeta("KickUserFromRoom.First")
		return
	}

	// Check if requester is the owner of the room.
	err = routes.AssertUserIsOwner(c, room.OwnerID)
	if err != nil {
		return
	}

	// Owner can't kick himself.
	if room.OwnerID == user.ID {
		c.Error(err).SetMeta("KickUserFromRoom.AssertUser")
		c.AbortWithError(http.StatusBadRequest, e.OwnerCantKickHimself{}).SetMeta("KickUserFromRoom.AssertUser")
		return
	}

	// Remove user from the connected users list of the room.
	err = room.RemoveUser(ctx, db, &user)
	if err != nil {
		c.Error(err).SetMeta("KickUserFromRoom.RemoveUser")
		c.AbortWithError(http.StatusInternalServerError, e.RoomNotModified{}).SetMeta("KickUserFromRoom.RemoveUser")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
