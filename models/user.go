package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	Name       string `json:"name" validate:"gte=2,lte=20,alphanumunicode"`
}

type UserUpdate struct {
	Name string `json:"name,omitempty" validate:"gte=2,lte=20,alphanumunicode"`
}

// ToUser converts a UserUpdate to a User
func (u *UserUpdate) ToUser() *User {
	return &User{
		Name: u.Name,
	}
}
