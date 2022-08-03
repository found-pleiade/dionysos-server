package utils

import (
	"log"

	"github.com/Brawdunoir/dionysos-server/models"
	"gorm.io/gorm"
)

// DeleteRoom deletes a room in the database
func DeleteRoom(id uint, db *gorm.DB) bool {
	result := db.Delete(&models.Room{}, id)

	if result.Error != nil {
		log.Printf("Failed to delete document: %v", result.Error)
		return false
	} else if result.RowsAffected < 1 {
		log.Printf("Failed to find document: %v", result.Error)
		return false
	}
	return true
}
