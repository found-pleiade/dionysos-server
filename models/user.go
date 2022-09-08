package models

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64       `gorm:"primarykey"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
	Name      string       `json:"name" binding:"required,gte=2,lte=20" example:"Diablox9"`
	Password  string       `json:"-"`
}

type UserUpdate struct {
	Name string `json:"name,omitempty" binding:"gte=2,lte=20" example:"Diablox9"`
}

// ToUser converts a UserUpdate to a User
func (u *UserUpdate) ToUser() *User {
	return &User{
		Name: u.Name,
	}
}

func (u *User) GetRoomID(ctx context.Context, db *gorm.DB) (uint64, error) {
	var roomID uint64
	err := db.WithContext(ctx).Select("room_id").Table("room_users").Where("user_id = ?", u.ID).Row().Scan(&roomID)
	if err != nil {
		return 0, err
	}
	return roomID, nil
}
