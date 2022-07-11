package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required"`
}

type UserUpdate struct {
	Username string `json:"username,omitempty"`
}
