package repository

import (
	"fmt"
	"github.com/La002/personal-crm/config"
	"github.com/La002/personal-crm/internal/entity"
	"github.com/La002/personal-crm/pkg/logger"
	"github.com/La002/personal-crm/pkg/postgres"
	"gorm.io/gorm"
)

type ContactRepo struct {
	DB *gorm.DB
}

func NewContactRepo(config *config.Configuration, l *logger.Logger) *ContactRepo {
	db := postgres.ConnectDB(config, l)
	return &ContactRepo{
		DB: db,
	}
}

func (r *ContactRepo) CreateContact(contact *entity.Contact) error {
	if err := r.DB.Create(contact).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContactRepo) SaveContact(newContact *entity.Contact) error {

	if err := r.DB.Save(newContact).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContactRepo) GetContact(id string) (entity.Contact, error) {
	var contact entity.Contact
	if err := r.DB.First(&contact, id).Error; err != nil {
		return entity.Contact{}, err
	}
	return contact, nil
}

func (r *ContactRepo) GetContactByEmail(email string) (entity.Contact, error) {
	var contact entity.Contact
	if err := r.DB.Where("email = ?", email).First(&contact).Error; err != nil {
		return entity.Contact{}, err
	}
	return contact, nil
}

func (r *ContactRepo) GetAllContacts() ([]entity.Contact, error) {
	var contacts []entity.Contact
	if err := r.DB.Select("id", "name", "relationship", "industry", "company", "birthday", "VIP", "last_met", "last_contacted", "last_update").Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *ContactRepo) SortContacts(sortBy string, sortOrder string) ([]entity.Contact, error) {
	var contacts []entity.Contact
	if err := r.DB.Order(sortBy + " " + sortOrder).Find(&contacts).Error; err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *ContactRepo) FilterContactsByVIPStatus() ([]entity.Contact, error) {
	var contacts []entity.Contact
	if err := r.DB.Where("vip = ?", true).Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *ContactRepo) DeleteContact(id string) error {
	result := r.DB.Where("id = ?", id).Delete(&entity.Contact{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no contact found with id %s", id)
	}
	return nil
}

func (r *ContactRepo) SearchContactByName(name string) ([]entity.Contact, error) {
	var contacts []entity.Contact
	if err := r.DB.Where("name ILIKE ?", "%"+name+"%").Find(&contacts).Error; err != nil {
		return nil, err
	}

	return contacts, nil
}
