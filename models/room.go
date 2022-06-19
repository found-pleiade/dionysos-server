package models

type Room struct {
	Name string `json:"name" binding:"required"`
}

type RoomUpdate struct {
	Name string `json:"name,omitempty"`
}
