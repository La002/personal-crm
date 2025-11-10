package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Relation string

const (
	Friend    Relation = "Friend"
	Family    Relation = "Family"
	Colleague Relation = "Colleague"
	School    Relation = "School"
	Network   Relation = "Network"
	Services  Relation = "Services"
)

type Contact struct {
	gorm.Model
	Name         string   `json:"name" gorm:"varchar(255);not null"`
	Relationship Relation `json:"relationship" gorm:"type:relation"`
	Industry     string   `json:"industry"`
	Company      string   `json:"company"`
	Birthday     string   `json:"birthday"`
	VIP          bool     `json:"vip"`
	DetailInfo
	VIPInfo
}

type DetailInfo struct {
	FamilyDetails
	ContactInfo
	Notes string `json:"notes"`
}

type FamilyDetails struct {
	Spouse   string `json:"spouse"`
	Children string `json:"children"`
}

type ContactInfo struct {
	Location    string `json:"location"`
	PhoneNumber string `json:"phone_number" gorm:"varchar(255);not null"`
	Email       string `json:"email" gorm:"varchar(255);not null"`
	LinkedIn    string `json:"linked_in"`
	Instagram   string `json:"instagram"`
	X           string `json:"x"`
}

type VIPInfo struct {
	LastMet       string `json:"last_met"`
	LastContacted string `json:"last_contacted"`
	LastUpdate    string `json:"last_update"`
	Status        string `json:"status"`
}

type DetailChanges struct {
	Id           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	PersonId     uint      `json:"person_id"`
	ChangedField string    `json:"changed_field"`
	OldValue     string    `json:"old_value"`
	NewValue     string    `json:"new_value"`
}
