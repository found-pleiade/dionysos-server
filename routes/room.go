//nolint:typecheck
package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/arangodb/go-driver"
	"github.com/gin-gonic/gin"
)

// CreateRoom creates a room in the aganro database
func CreateRoom(c *gin.Context) {
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	col, err := db.Collection(ctx, database.RoomsCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot access database collection"})
		log.Printf("Failed to access collection: %v", err)
		return
	}

	meta, err := col.CreateDocument(ctx, room)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Room not created"})
		log.Printf("Failed to create document: %v", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": meta.Key})
}

// GetRoom returns a room from the aganro database
func GetRoom(c *gin.Context) {
	var result models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id := c.Param("id")

	col, err := db.Collection(ctx, database.RoomsCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot access database collection"})
		log.Printf("Failed to access collection: %v", err)
		return
	}

	_, err = col.ReadDocument(ctx, id, &result)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Room not found"})
		log.Printf("Failed to get document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": result})
}

// UpdateRoom updates a room in the aganro database
func UpdateRoom(c *gin.Context) {
	var roomUpdate models.RoomUpdate
	var patchedRoom models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id := c.Param("id")

	if err := c.ShouldBindJSON(&roomUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	col, err := db.Collection(ctx, database.RoomsCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot access database collection"})
		log.Printf("Failed to access collection: %v", err)
		return
	}

	_, err = col.UpdateDocument(driver.WithReturnNew(ctx, &patchedRoom), id, roomUpdate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Room not modified"})
		log.Printf("Failed to modify document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": patchedRoom})
}

// DeleteRoom deletes a room in the aganro database
func DeleteRoom(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id := c.Param("id")

	col, err := db.Collection(ctx, database.RoomsCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot access database collection"})
		log.Printf("Failed to access collection: %v", err)
		return
	}

	_, err = col.RemoveDocument(ctx, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Room not deleted"})
		log.Printf("Failed to delete document: %v", err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
