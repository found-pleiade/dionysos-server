package models

type Room struct {
	Name string `json:"name" binding:"required"`
}
