//nolint:typecheck
package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/gin-gonic/gin"
)

// CreateUser creates a user in the aganro database
func CreateUser(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 500*time.Millisecond)
	defer cancelCtx()
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	col, err := db.Collection(ctx, database.UsersCollection)
	if err != nil {
		fmt.Println(err)
	}

	meta, err := col.CreateDocument(ctx, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not created"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created", "id": meta.Key})
}

// GetUser returns a user from the aganro database
func GetUser(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 500*time.Millisecond)
	defer cancelCtx()
	var result models.User
	id := c.Param("id")

	col, err := db.Collection(ctx, database.UsersCollection)
	if err != nil {
		fmt.Println(err)
	}

	_, err = col.ReadDocument(ctx, id, &result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User found", "user": result})
}

// UpdateUser updates a user in the aganro database
func UpdateUser(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 500*time.Millisecond)
	defer cancelCtx()
	var user models.User
	id := c.Param("id")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	col, err := db.Collection(ctx, database.UsersCollection)
	if err != nil {
		fmt.Println(err)
	}

	patch := map[string]interface{}{
		"username": user.Username,
	}

	meta, err := col.UpdateDocument(ctx, id, patch)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not modified"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User modified", "id": meta.Key})
}

// DeleteUser deletes a user in the aganro database
func DeleteUser(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 500*time.Millisecond)
	defer cancelCtx()
	id := c.Param("id")

	col, err := db.Collection(ctx, database.UsersCollection)
	if err != nil {
		fmt.Println(err)
	}

	meta, err := col.RemoveDocument(ctx, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted", "id": meta.Key})
}
