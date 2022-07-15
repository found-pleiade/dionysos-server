package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model `json:"-"`
	Name       string `json:"name" validate:"required,gte=2,lte=20,alphanumunicode"`
}

type RoomUpdate struct {
	Name string `json:"name,omitempty" validate:"gte=2,lte=20,alphanumunicode"`
}

// ToRoom converts a RoomUpdate to a Room
func (ru *RoomUpdate) ToRoom() *Room {
	return &Room{
		Name: ru.Name,
	}
}
