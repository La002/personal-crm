package entity

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	UserID                uint   `gorm:"not null;index"`
	ContactID             uint   `gorm:"not null;index"`
	Title                 string `gorm:"not null"`
	EventDate             string `gorm:"not null"` // Date of the event (YYYY-MM-DD)
	Recurrence            string // "none", "monthly", "yearly"
	GoogleCalendarEventID string
}
