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

// CreateUser creates a user in the database
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

// GetUser returns a user from the database
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

// UpdateUser updates a user in the database
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

// DeleteUser deletes a user in the database
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
