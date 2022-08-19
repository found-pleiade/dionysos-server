package models

import (
	"database/sql"
	"time"
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
	Name string `json:"name,omitempty" binding:"gte=2,lte=20"`
}

// ToUser converts a UserUpdate to a User
func (u *UserUpdate) ToUser() *User {
	return &User{
		Name: u.Name,
	}
}
