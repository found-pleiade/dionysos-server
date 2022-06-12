//nolint:typecheck
package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/gin-gonic/gin"
)

// CreateUser creates a user in the aganro database
func CreateUser(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()
	var user models.User
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
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()
	var result models.User
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
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()
	var user models.User
	id := c.Param("id")

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

	patch := map[string]interface{}{
		"username": user.Username,
	}

	meta, err := col.UpdateDocument(ctx, id, patch)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not modified"})
		log.Printf("Failed to modify document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": meta.Key})
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

	meta, err := col.RemoveDocument(ctx, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not deleted"})
		log.Printf("Failed to delete document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": meta.Key})
}
