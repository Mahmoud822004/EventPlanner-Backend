package models

import (
	"time"
)

type Event struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	Time        string    `json:"time"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	OrganizerID uint      `json:"organizer_id"`
	Invitations []Invitation  `gorm:"foreignKey:EventID"`
}

type Invitation struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	EventID uint   `json:"event_id"`
	UserID  uint   `json:"user_id"`
	Role    string `json:"role"` // "organizer" or "attendee"
	Status  string `json:"status"` // going, maybe, not_going

	User User `gorm:"foreignKey:UserID" json:"user"`
}