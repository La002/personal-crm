package repository

import "github.com/La002/personal-crm/internal/entity"

type ContactDao interface {
	CreateContact(Contact *entity.Contact) error
	SaveContact(Contact *entity.Contact) error
	GetContact(id string) (entity.Contact, error)
	GetContactByEmail(email string) (entity.Contact, error)
	GetAllContacts() ([]entity.Contact, error)
	SortContacts(sortBy string, sortOrder string) ([]entity.Contact, error)
	FilterContactsByVIPStatus() ([]entity.Contact, error)
	DeleteContact(id string) error
	SearchContactByName(name string) ([]entity.Contact, error)
}
