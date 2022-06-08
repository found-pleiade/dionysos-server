//nolint:typecheck
package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/gin-gonic/gin"
)

// CreateRoom creates a room in the aganro database
func CreateRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	col, err := db.Collection(context.Background(), database.RoomsCollection)
	if err != nil {
		fmt.Println(err)
	}

	meta, err := col.CreateDocument(context.Background(), room)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Room not created"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room created", "id": meta.Key})
}

// GetRoom returns a room from the aganro database
func GetRoom(c *gin.Context) {
	var result models.Room
	id := c.Param("id")

	col, err := db.Collection(context.Background(), database.RoomsCollection)
	if err != nil {
		fmt.Println(err)
	}

	_, err = col.ReadDocument(context.Background(), id, &result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room found", "room": result})
}

// UpdateRoom updates a room in the aganro database
func UpdateRoom(c *gin.Context) {
	var room models.Room
	id := c.Param("id")

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	col, err := db.Collection(context.Background(), database.RoomsCollection)
	if err != nil {
		fmt.Println(err)
	}

	patch := map[string]interface{}{
		"name": room.Name,
	}

	meta, err := col.UpdateDocument(context.Background(), id, patch)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Room not modified"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room modified", "id": meta.Key})
}

// DeleteRoom deletes a room in the aganro database
func DeleteRoom(c *gin.Context) {
	id := c.Param("id")

	col, err := db.Collection(context.Background(), database.RoomsCollection)
	if err != nil {
		fmt.Println(err)
	}

	meta, err := col.RemoveDocument(context.Background(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Room not deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room deleted", "id": meta.Key})
}
