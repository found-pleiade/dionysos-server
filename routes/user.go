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
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUser creates a user in the database
func CreateUser(c *gin.Context) {
	var user models.User

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to validate: %v", err)
		return
	}

	err := db.WithContext(ctx).Create(&user).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not created"})
		log.Printf("Failed to create document: %v", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"uri": "/users/" + fmt.Sprint(user.ID)})
}

// GetUser returns a user from the database
func GetUser(c *gin.Context) {
	var user models.User

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		log.Printf("Failed to convert user ID: %v", err)
	}

	err = db.WithContext(ctx).First(&user, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		log.Printf("Failed to find document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser updates a user in the database
func UpdateUser(c *gin.Context) {
	var userUpdate models.UserUpdate
	var patchedUser models.User

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		log.Printf("Failed to convert user ID: %v", err)
	}

	// Test if data is valid
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}
	if err := validate.Struct(userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to validate: %v", err)
		return
	}

	err = db.WithContext(ctx).First(&patchedUser, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			log.Printf("Failed to find document: %v", err)
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not updated"})
			log.Printf("Failed to modify document: %v", err)
			return
		}
	}

	err = db.WithContext(ctx).Model(&patchedUser).Updates(userUpdate.ToUser()).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not modified"})
		log.Printf("Failed to modify document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": patchedUser})
}

// DeleteUser deletes a user in the database
func DeleteUser(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		log.Printf("Failed to convert user ID: %v", err)
	}

	result := db.WithContext(ctx).Delete(&models.User{}, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not deleted"})
		log.Printf("Failed to delete document: %v", result.Error)
		return
	} else if result.RowsAffected < 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		log.Printf("Failed to find document: %v", result.Error)
		return
	}

	c.JSON(http.StatusOK, nil)
}
