package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type Room struct {
	ID        uint64       `gorm:"primaryKey;autoincrement:false" json:"-"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
	Name      string       `json:"name" binding:"required,gte=2,lte=20" example:"BirthdayParty"`
	OwnerID   uint64       `json:"ownerID"`
	Users     []User       `json:"users" gorm:"many2many:room_users"`
	MediaURL  string       `json:"mediaURL"`
	Pause     bool         `json:"pause"`
	Timestamp uint64       `json:"timestamp"`
}

type RoomUpdate struct {
	Name     string `json:"name,omitempty" binding:"gte=2,lte=20" example:"BirthdayParty"`
	MediaURL string `json:"mediaURL,omitempty" binding:"gt=0" example:"https://www.youtube.com/watch?v=dQw4w9WgXcQ"`
}

// ToRoom converts a RoomUpdate to a Room
func (ru *RoomUpdate) ToRoom() *Room {
	return &Room{
		Name:     ru.Name,
		MediaURL: ru.MediaURL,
	}
}

// GetRoom gets a room by its ID and sets its Users field before returning it.
func (r *Room) GetRoom(ctx context.Context, db *gorm.DB, id uint64) error {
	err := db.WithContext(ctx).First(&r, id).Error
	if err != nil {
		return err
	}
	err = db.WithContext(ctx).Model(&r).Association("Users").Find(&r.Users)
	if err != nil {
		return err
	}

	return err
}

// RemoveUser removes a user from a room.
func (r *Room) RemoveUser(ctx context.Context, db *gorm.DB, user *User) error {
	if !slices.Contains(r.Users, *user) {
		return errors.New("User not connected to room")
	}

	return db.WithContext(ctx).Model(&r).Association("Users").Delete(user)
}
