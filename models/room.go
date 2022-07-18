package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model `json:"-"`
	Name       string `json:"name" binding:"required,gte=2,lte=20"`
}

type RoomUpdate struct {
	Name string `json:"name,omitempty" binding:"gte=2,lte=20"`
}

// ToRoom converts a RoomUpdate to a Room
func (ru *RoomUpdate) ToRoom() *Room {
	return &Room{
		Name: ru.Name,
	}
}
