package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	GoogleId string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Picture  string

	AccessToken  string `gorm:"type:text"`
	RefreshToken string `gorm:"type:text"`
	TokenExpiry  time.Time

	CalendarSyncEnabled bool
	LastCalendarSync    time.Time
}
