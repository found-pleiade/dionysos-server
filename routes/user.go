//nolint:typecheck
package routes

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/Brawdunoir/dionysos-server/utils"
	e "github.com/Brawdunoir/dionysos-server/utils/errors"
	l "github.com/Brawdunoir/dionysos-server/utils/logger"
	routes "github.com/Brawdunoir/dionysos-server/utils/routes"
	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary      Creates a user. Needed for further request during authentication.
// @Description  Creates a user. You will need to use BasicAuth to authenticate with the created user, using its ID and password produced by this endpoint.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user body models.UserUpdate true "User object"
// @Success      201	{object} utils.CreateResponse "User created"
// @Failure      400	{object} utils.ErrorResponse "Invalid request"
// @Failure      500	{object} utils.ErrorResponse "Internal server error"
// @Router       /users [post]
func CreateUser(c *gin.Context) {
	var u models.UserUpdate
	rand.Seed(time.Now().UnixNano())
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&u); err != nil {
		c.Error(err).SetMeta("CreateUser.ShouldBindJSON")
		c.AbortWithError(http.StatusBadRequest, e.FailJSONBind{}).SetMeta("CreateUser.ShouldBindJSON")
		return
	}

	// Create 32 bytes random password
	password := make([]byte, 32)
	rand.Read(password)
	user := u.ToUser()

	user.Password = fmt.Sprintf("%x", password)

	err := db.WithContext(ctx).Create(&user).Error

	if err != nil {
		c.Error(err).SetMeta("CreateUser.Create")
		c.AbortWithError(http.StatusInternalServerError, e.UserNotCreated{}).SetMeta("CreateUser.Create")
		return
	}

	c.JSON(http.StatusCreated, routes.CreateResponse{URI: "/users/" + fmt.Sprint(user.ID), Password: user.Password})
}

// GetUser godoc
// @Summary      Gets a user.
// @Tags         Users
// @Security     BasicAuth
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 	{object} models.User
// @Failure      400 	{object} utils.ErrorResponse "Invalid request"
// @Failure      401 	{object} utils.ErrorResponse "User not authorized"
// @Failure      404 	{object} utils.ErrorResponse "User not found"
// @Router       /users/{id} [get]
func GetUser(c *gin.Context) {
	var user models.User
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(err).SetMeta("GetUser.ParseUint")
		c.AbortWithError(http.StatusBadRequest, e.InvalidID{}).SetMeta("GetUser.ParseUint")
		return
	}

	err = db.WithContext(ctx).First(&user, id).Error

	if err != nil {
		c.Error(err).SetMeta("GetUser.First")
		c.AbortWithError(http.StatusNotFound, e.UserNotFound{}).SetMeta("GetUser.First")
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary      Updates a user.
// @Tags         Users
// @Security     BasicAuth
// @Accept       json
// @Produce      json
// @Param        id   path int               true "User ID"
// @Param        user body models.UserUpdate true "User object"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "User not found"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	var u models.UserUpdate
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	patchedUser, err := routes.ExtractUserFromContext(c)
	if err != nil {
		c.Error(err).SetMeta("UpdateUser.ExtractUserFromContext")
		c.AbortWithError(http.StatusInternalServerError, e.UserNotInContext{}).SetMeta("UpdateUser.ExtractUserFromContext")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(err).SetMeta("UpdateUser.ParseUint")
		c.AbortWithError(http.StatusBadRequest, e.InvalidID{}).SetMeta("UpdateUser.ParseUint")
		return
	}

	// Assert the request is coming from the right user.
	if err := routes.AssertUserIsOwner(c, id); err != nil {
		return
	}

	// Test if data is valid.
	if err := c.ShouldBindJSON(&u); err != nil {
		c.Error(err).SetMeta("UpdateUser.BindJSON")
		c.AbortWithError(http.StatusBadRequest, e.FailJSONBind{}).SetMeta("UpdateUser.BindJSON")
		return
	}

	err = db.WithContext(ctx).Model(&patchedUser).Updates(u.ToUser()).Error
	if err != nil {
		c.Error(err).SetMeta("UpdateUser.Updates")
		c.AbortWithError(http.StatusInternalServerError, e.UserNotModified{}).SetMeta("UpdateUser.Updates")
		return
	}

	// If the user has a room, broadcast its rename to room members.
	roomID, err := patchedUser.GetRoomID(ctx, db)
	if err == nil {
		stream, err := utils.GetStream(roomID, roomStreamsList)
		if err != nil {
			l.Logger.Infof("Failed to get stream: %v", err)
		} else {
			stream.Distribute(SSEMessage)
		}
	}

	c.JSON(http.StatusNoContent, nil)
}

// DeleteUser godoc
// @Summary      Deletes a user. Should be used when disconnecting a user.
// @Tags         Users
// @Security     BasicAuth
// @Param        id path int true "User ID"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "User not found"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(err).SetMeta("DeleteUser.ParseUint")
		c.AbortWithError(http.StatusBadRequest, e.InvalidID{}).SetMeta("DeleteUser.ParseUint")
	}

	// Assert the request is coming from the right user.
	if err := routes.AssertUserIsOwner(c, id); err != nil {
		return
	}

	result := db.WithContext(ctx).Delete(&models.User{}, id)

	if result.Error != nil {
		c.Error(result.Error).SetMeta("DeleteUser.Delete.NotDeleted")
		c.AbortWithError(http.StatusInternalServerError, e.UserNotDeleted{}).SetMeta("DeleteUser.Delete.NotDeleted")
		return
	} else if result.RowsAffected < 1 {
		c.Error(result.Error).SetMeta("DeleteUser.Delete.NotFound")
		c.AbortWithError(http.StatusNotFound, e.UserNotFound{}).SetMeta("DeleteUser.Delete.NotFound")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
