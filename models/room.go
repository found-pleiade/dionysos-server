package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Room struct {
	ID        uint          `gorm:"primarykey" json:"-"`
	CreatedAt time.Time     `json:"-"`
	UpdatedAt time.Time     `json:"-"`
	DeletedAt sql.NullTime  `gorm:"index" json:"-"`
	Name      string        `json:"name" binding:"required,gte=2,lte=20" example:"BirthdayParty"`
	OwnerID   uint          `json:"ownerID,omitempty"`
	UsersID   pq.Int64Array `json:"usersID,omitempty" gorm:"type:integer[]"`
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
