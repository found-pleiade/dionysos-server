package models

import (
	"context"
	"database/sql"
	"time"

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
}

type RoomUpdate struct {
	Name string `json:"name,omitempty" binding:"gte=2,lte=20"`
}

// ToRoom converts a RoomUpdate to a Room
func (ru *RoomUpdate) ToRoom() *Room {
	return &Room{
		Name: ru.Name,
	}
}

func (r *Room) GetRoom(ctx context.Context, db *gorm.DB, id uint64) error {
	err := db.WithContext(ctx).First(&r, id).Error
	if err != nil {
		return err
	}
	r.Users, err = r.GetAllUsers(ctx, db)
	if err != nil {
		return err
	}

	return err
}

// Retrieve user list with edger loading languages
func (r *Room) GetAllUsers(ctx context.Context, db *gorm.DB) ([]User, error) {
	var users []User
	err := db.WithContext(ctx).Model(&r).Association("Users").Find(&users)
	return users, err
}

func (r *Room) RemoveUser(ctx context.Context, db *gorm.DB, user *User) error {
	return db.WithContext(ctx).Model(&r).Association("Users").Delete(user)
}
