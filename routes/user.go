//nolint:typecheck
package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	driver "github.com/arangodb/go-driver"

	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/gin-gonic/gin"
)

// CreateUser creates a user in the aganro database
func CreateUser(c *gin.Context) {
	var user models.User

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	col, err := db.Collection(ctx, database.UsersCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot access database collection"})
		log.Printf("Failed to access collection: %v", err)
		return
	}

	meta, err := col.CreateDocument(ctx, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not created"})
		log.Printf("Failed to create document: %v", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": meta.Key})
}

// GetUser returns a user from the aganro database
func GetUser(c *gin.Context) {
	var result models.User

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id := c.Param("id")

	col, err := db.Collection(ctx, database.UsersCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot access database collection"})
		log.Printf("Failed to access collection: %v", err)
		return
	}

	_, err = col.ReadDocument(ctx, id, &result)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		log.Printf("Failed to find document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": result})
}

// UpdateUser updates a user in the aganro database
func UpdateUser(c *gin.Context) {
	var userUpdate models.UserUpdate
	var patchedUser models.User

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id := c.Param("id")

	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	if isNil := userUpdate == (models.UserUpdate{}); isNil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data to update"})
		log.Printf("Failed to bind JSON: No data to update")
		return
	}

	col, err := db.Collection(ctx, database.UsersCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot access database collection"})
		log.Printf("Failed to access collection: %v", err)
		return
	}

	_, err = col.UpdateDocument(driver.WithReturnNew(ctx, &patchedUser), id, userUpdate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not modified"})
		log.Printf("Failed to modify document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": patchedUser})
}

// DeleteUser deletes a user in the aganro database
func DeleteUser(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id := c.Param("id")

	col, err := db.Collection(ctx, database.UsersCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot access database collection"})
		log.Printf("Failed to access collection: %v", err)
		return
	}

	_, err = col.RemoveDocument(ctx, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not deleted"})
		log.Printf("Failed to delete document: %v", err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
