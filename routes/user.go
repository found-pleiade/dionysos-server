//nolint:typecheck
package routes

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Brawdunoir/dionysos-server/models"
	routes "github.com/Brawdunoir/dionysos-server/utils/routes"
	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary      Creates a user. Needed for further request during authentication.
// @Description  Creates a user. You will need to use BasicAuth to authenticate with the created user, using its ID and password produced by this endpoint.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param user body models.UserUpdate true "User object"
// @Success      201 {object} utils.CreateResponse "User created"
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /users [post]
func CreateUser(c *gin.Context) {
	var u models.UserUpdate
	rand.Seed(time.Now().UnixNano())

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse(err.Error()))
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	// Create 32 bytes random password
	password := make([]byte, 32)
	rand.Read(password)
	user := u.ToUser()

	user.Password = fmt.Sprintf("%x", password)

	err := db.WithContext(ctx).Create(&user).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("User not created"))
		log.Printf("Failed to create document: %v", err)
		return
	}

	c.JSON(http.StatusCreated, routes.CreateResponse{URI: "/users/" + fmt.Sprint(user.ID), Password: user.Password})
}

// GetUser godoc
// @Summary      Gets a user.
// @Tags         Users
// @Security BasicAuth
// @Produce      json
// @Param 			 id path int true "User ID"
// @Success      200 {object} models.User
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "User not found"
// @Router       /users/{id} [get]
func GetUser(c *gin.Context) {
	var user models.User

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse("Invalid user ID"))
		log.Printf("Failed to convert user ID: %v", err)
	}

	err = db.WithContext(ctx).First(&user, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, routes.CreateErrorResponse("User not found"))
		log.Printf("Failed to find document: %v", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary      Updates a user.
// @Tags         Users
// @Security BasicAuth
// @Accept       json
// @Produce      json
// @Param id path int true "User ID"
// @Param user body models.UserUpdate true "User object"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      401 {object} utils.ErrorResponse "User not authorized"
// @Failure      404 {object} utils.ErrorResponse "User not found"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	var u models.UserUpdate
	patchedUser, err := routes.ExtractUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("User not found in context. Has it been set in the middleware?"))
		log.Printf("Failed to extract user from context: %v", err)
	}

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse("Invalid user ID"))
		log.Printf("Failed to convert user ID: %v", err)
	}

	// Assert the request is coming from the right user.
	if err := routes.AssertUser(c, id); err != nil {
		log.Printf("Failed to assert user: %v", err)
		return
	}

	// Test if data is valid.
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse(err.Error()))
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	err = db.WithContext(ctx).Model(&patchedUser).Updates(u.ToUser()).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("User not modified"))
		log.Printf("Failed to modify document: %v", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// DeleteUser godoc
// @Summary      Deletes a user. Should be used when disconnecting a user.
// @Tags         Users
// @Security BasicAuth
// @Param id path int true "User ID"
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
		c.JSON(http.StatusBadRequest, routes.CreateErrorResponse("Invalid user ID"))
		log.Printf("Failed to convert user ID: %v", err)
	}

	// Assert the request is coming from the right user.
	if err := routes.AssertUser(c, id); err != nil {
		log.Printf("Failed to assert user: %v", err)
		return
	}

	result := db.WithContext(ctx).Delete(&models.User{}, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, routes.CreateErrorResponse("User not deleted"))
		log.Printf("Failed to delete document: %v", result.Error)
		return
	} else if result.RowsAffected < 1 {
		c.JSON(http.StatusNotFound, routes.CreateErrorResponse("User not found"))
		log.Printf("Failed to find document: %v", result.Error)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
