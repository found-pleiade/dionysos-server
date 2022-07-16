package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        uint         `gorm:"primarykey" json:"-"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
	Name      string       `json:"name" binding:"required" example:"Diablox9"`
}

type UserUpdate struct {
	Name string `json:"name,omitempty"`
}

// ToUser converts a UserUpdate to a User
func (u *UserUpdate) ToUser() *User {
	return &User{
		Name: u.Name,
	}
}
