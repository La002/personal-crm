package repository

import (
	"time"

	"github.com/La002/personal-crm/pkg/entity"
)

type ContactDao interface {
	CreateContact(Contact *entity.Contact) error
	SaveContact(Contact *entity.Contact) error
	GetContact(id string, userID uint) (entity.Contact, error)
	GetContactByEmail(email string, userID uint) (entity.Contact, error)
	GetAllContacts(userID uint) ([]entity.Contact, error)
	SortContacts(sortBy string, sortOrder string, userID uint) ([]entity.Contact, error)
	DeleteContact(id string, userID uint) error
	FilterContactsByVipStatus(flag bool, userID uint) ([]entity.Contact, error)
	SearchContactsByRelationship(filter string, userID uint) ([]entity.Contact, error)
	UpdateCalendarSync(contactID string, userID uint, eventID string, synced bool) error

	// Event methods
	CreateEvent(event *entity.Event) error
	GetEventsByContact(contactID, userID uint) ([]entity.Event, error)
	GetEventByID(eventID, userID uint) (entity.Event, error)
	DeleteEvent(eventID, userID uint) error
	UpdateEventGoogleID(eventID, userID uint, googleEventID string) error
	GetUpcomingEvents(userID uint, days int) ([]entity.Event, error)

	// MCP specific
	GetAllContactsWithLimit(userID uint, limit int) ([]entity.Contact, error)
	SearchContactsAdvanced(userID uint, filters ContactSearchFilters) ([]entity.Contact, error)
	UpdateContactFields(id string, userID uint, updates map[string]interface{}) error
	AppendNotes(id string, userID uint, newNotes string) error
}

type UserDao interface {
	CreateUser(user *entity.User) error
	GetUserByGoogleID(id string) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
	GetUserByID(id uint) (entity.User, error)
	UpdateUser(user *entity.User) error
}

type ContactSearchFilters struct {
	Relationship        *string
	VipOnly             *bool
	Location            *string
	LastContactedBefore *time.Time
	Limit               int
}
