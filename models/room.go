package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model `json:"-"`
	Name       string `json:"name" binding:"required"`
}

type RoomUpdate struct {
	Name string `json:"name,omitempty"`
}

// ToRoom converts a RoomUpdate to a Room
func (ru *RoomUpdate) ToRoom() *Room {
	return &Room{
		Name: ru.Name,
	}
}
