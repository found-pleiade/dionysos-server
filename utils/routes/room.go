package utils

import (
	"log"

	"github.com/Brawdunoir/dionysos-server/models"
	"gorm.io/gorm"
)

// DeleteRoom deletes a room in the database
func DeleteRoom(id uint, db *gorm.DB) error {
	result := db.Delete(&models.Room{}, id)

	if result.Error != nil {
		log.Printf("Failed to delete document: %v", result.Error)
		return result.Error
	} else if result.RowsAffected < 1 {
		log.Printf("Failed to find document: %v", result.Error)
		return result.Error
	}
	return nil
}

// ChangeOwner changes the owner of a room in the database
func ChangeOwner(roomUpdate models.RoomUpdate, db *gorm.DB) error {
	var room models.Room

	// We already removed the user from the slice, so we can get the first element as new Owner
	roomUpdate.OwnerID = roomUpdate.UsersID[0]
	err := db.Model(&room).Updates(roomUpdate.ToRoom()).Error
	if err != nil {
		return err
	}
	return nil
}
