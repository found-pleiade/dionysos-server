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
	c.BindJSON(&user)

	col, err := db.Collection(context.TODO(), database.UsersCollection)
	if err != nil {
		fmt.Println(err)
	}

	meta, err := col.CreateDocument(context.TODO(), user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not created"})
		log.Fatalf("Failed to create documents: %v", err)
	}

	fmt.Printf("Created document with keys '%s' in collection '%s' in database '%s'\n", meta.Key, col.Name(), db.Name())

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
		log.Fatalf("Failed to read documents: %v", err)
	}

	fmt.Printf("Read document with keys '%s' in collection '%s' in database '%s'\n", meta.Key, col.Name(), db.Name())

	c.JSON(http.StatusOK, gin.H{"message": "User found", "user": result})
}
