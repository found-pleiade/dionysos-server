package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name string `json:"name" binding:"required"`
}

type RoomUpdate struct {
	gorm.Model
	Name string `json:"name,omitempty"`
}
