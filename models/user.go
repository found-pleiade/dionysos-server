package models

type User struct {
	Username string `json:"username" binding:"required"`
}

type UserUpdate struct {
	Username string `json:"username,omitempty"`
}
