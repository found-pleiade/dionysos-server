package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	Name       string `json:"name" binding:"required,gte=2,lte=20"`
	Password   string `json:"-"`
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
