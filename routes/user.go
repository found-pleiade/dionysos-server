//nolint:typecheck
package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/gin-gonic/gin"
)

// CreateUser creates a user in the aganro database
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	col, err := db.Collection(context.TODO(), database.UsersCollection)
	if err != nil {
		fmt.Println(err)
	}

	meta, err := col.CreateDocument(context.TODO(), user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not created"})
		log.Printf("Failed to create documents: %v", err)
		return
	}

	fmt.Printf("Created document with key '%s' in collection '%s' in database '%s'\n", meta.Key, col.Name(), db.Name())

	c.JSON(http.StatusOK, gin.H{"message": "User created", "id": meta.Key})
}

// GetUser returns a user from the aganro database
func GetUser(c *gin.Context) {
	id := c.Param("id")

	col, err := db.Collection(context.TODO(), database.UsersCollection)
	if err != nil {
		fmt.Println(err)
	}

	var result models.User
	meta, err := col.ReadDocument(context.TODO(), id, &result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
		log.Printf("Failed to read documents: %v", err)
		return
	}

	fmt.Printf("Read document with key '%s' in collection '%s' in database '%s'\n", meta.Key, col.Name(), db.Name())

	c.JSON(http.StatusOK, gin.H{"message": "User found", "user": result})
}

// UpdateUser updates a user in the aganro database
func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	col, err := db.Collection(context.TODO(), database.UsersCollection)
	if err != nil {
		fmt.Println(err)
	}

	patch := map[string]interface{}{
		"username": user.Username,
	}

	meta, err := col.UpdateDocument(context.TODO(), id, patch)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not modified"})
		log.Printf("Failed to modify user: %v", err)
		return
	}

	fmt.Printf("Modified user with key '%s' in collection '%s' in database '%s'\n", meta.Key, col.Name(), db.Name())

	c.JSON(http.StatusOK, gin.H{"message": "User modified", "id": meta.Key})
}

// DeleteUser creates a user in the aganro database
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	col, err := db.Collection(context.TODO(), database.UsersCollection)
	if err != nil {
		fmt.Println(err)
	}

	meta, err := col.RemoveDocument(context.TODO(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not deleted"})
		log.Printf("Failed to delete user: %v", err)
	}

	fmt.Printf("Deleted user with key '%s' in collection '%s' in database '%s'\n", meta.Key, col.Name(), db.Name())

	c.JSON(http.StatusOK, gin.H{"message": "User deleted", "id": meta.Key})
}
