package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	Name       string `json:"name" binding:"required"`
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
