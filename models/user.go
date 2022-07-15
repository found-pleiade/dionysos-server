package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	Username   string `json:"username" binding:"required"`
}

type UserUpdate struct {
	Username string `json:"username,omitempty"`
}

// ToUser converts a UserUpdate to a User
func (u *UserUpdate) ToUser() *User {
	return &User{
		Username: u.Username,
	}
}
