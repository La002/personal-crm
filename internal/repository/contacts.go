package repository

import (
	"fmt"
	"time"

	"github.com/La002/personal-crm/config"
	"github.com/La002/personal-crm/internal/entity"
	"github.com/La002/personal-crm/pkg/logger"
	"github.com/La002/personal-crm/pkg/postgres"
	"gorm.io/gorm"
)

type ContactRepo struct {
	DB *gorm.DB
}

// MCP specific methods
func (r *ContactRepo) AppendNotes(id string, userID uint, newNotes string) error {
	var contact entity.Contact

	// First, get the current contact
	if err := r.DB.Where("id = ? AND user_id = ?", id, userID).First(&contact).Error; err != nil {
		return err
	}

	// Append new notes with timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	separator := "\n---\n"
	if contact.Notes == "" {
		contact.Notes = fmt.Sprintf("[%s]\n%s", timestamp, newNotes)
	} else {
		contact.Notes = fmt.Sprintf("%s%s[%s]\n%s", contact.Notes, separator, timestamp, newNotes)
	}

	// Update the notes field
	result := r.DB.Model(&entity.Contact{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("notes", contact.Notes)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ContactRepo) UpdateContactFields(id string, userID uint, updates map[string]interface{}) error {
	result := r.DB.Model(&entity.Contact{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no contact found with id %s", id)
	}

	return nil
}

func (r *ContactRepo) SearchContactsAdvanced(userID uint, filters ContactSearchFilters) ([]entity.Contact, error) {
	var contacts []entity.Contact
	query := r.DB.Where("user_id = ?", userID)

	// Apply relationship filter with prefix matching
	if filters.Relationship != nil && *filters.Relationship != "" {
		query = query.Where("relationship::TEXT ILIKE ?", *filters.Relationship+"%")
	}

	// Apply VIP filter
	if filters.VipOnly != nil {
		query = query.Where("vip = ?", *filters.VipOnly)
	}

	// Apply location filter with prefix matching
	if filters.Location != nil && *filters.Location != "" {
		query = query.Where("location ILIKE ?", *filters.Location+"%")
	}

	// Apply last contacted before filter
	if filters.LastContactedBefore != nil {
		query = query.Where("last_contacted < ?", *filters.LastContactedBefore)
	}

	// Apply limit
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}

	if err := query.Find(&contacts).Error; err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *ContactRepo) GetAllContactsWithLimit(userID uint, limit int) ([]entity.Contact, error) {
	var contacts []entity.Contact
	query := r.DB.Where("user_id = ?", userID)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
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

func (r *ContactRepo) GetContact(id string, userID uint) (entity.Contact, error) {
	var contact entity.Contact
	if err := r.DB.Where("id = ? AND user_id = ?", id, userID).First(&contact).Error; err != nil {
		return entity.Contact{}, err
	}
	return contact, nil
}

func (r *ContactRepo) GetContactByEmail(email string, userID uint) (entity.Contact, error) {
	var contact entity.Contact
	if err := r.DB.Where("email = ? AND user_id = ?", email, userID).First(&contact).Error; err != nil {
		return entity.Contact{}, err
	}
	return contact, nil
}

func (r *ContactRepo) GetAllContacts(userID uint) ([]entity.Contact, error) {
	var contacts []entity.Contact
	if err := r.DB.Where("user_id = ?", userID).Select("id", "name", "relationship", "industry", "company", "birthday", "Vip", "last_met", "last_contacted", "last_update").Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *ContactRepo) SortContacts(sortBy string, sortOrder string, userID uint) ([]entity.Contact, error) {
	var contacts []entity.Contact
	if err := r.DB.Where("user_id = ?", userID).Order(sortBy + " " + sortOrder).Find(&contacts).Error; err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *ContactRepo) FilterContactsByVipStatus(flag bool, userID uint) ([]entity.Contact, error) {
	var contacts []entity.Contact
	if err := r.DB.Where("vip = ? AND user_id = ?", flag, userID).Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *ContactRepo) DeleteContact(id string, userID uint) error {
	result := r.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Contact{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no contact found with id %s", id)
	}
	return nil
}

func (r *ContactRepo) SearchContactsByRelationship(relation string, userID uint) ([]entity.Contact, error) {
	var contacts []entity.Contact
	if err := r.DB.Where("relationship::TEXT ILIKE ? AND user_id = ?", "%"+relation+"%", userID).Find(&contacts).Error; err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *ContactRepo) UpdateCalendarSync(contactID string, userID uint, eventID string, synced bool) error {
	return r.DB.Model(&entity.Contact{}).
		Where("id = ? AND user_id = ?", contactID, userID).
		Updates(map[string]interface{}{
			"google_calendar_event_id": eventID,
			"calendar_sync_enabled":    synced,
			"calendar_synced_at":       time.Now(),
		}).Error
}

// Dashboard methods
func (r *ContactRepo) GetUpcomingBirthdays(userID uint, days int) ([]entity.Contact, error) {
	var contacts []entity.Contact

	// Get all contacts with birthdays
	err := r.DB.Where("user_id = ? AND birthday != ''", userID).Find(&contacts).Error
	if err != nil {
		return nil, err
	}

	// Filter for upcoming birthdays
	upcoming := []entity.Contact{}
	now := time.Now()

	for _, contact := range contacts {
		if contact.Birthday == "" {
			continue
		}

		// Parse birthday
		bday, err := time.Parse("2006-01-02", contact.Birthday)
		if err != nil {
			continue
		}

		// Calculate next birthday this year
		nextBday := time.Date(now.Year(), bday.Month(), bday.Day(), 0, 0, 0, 0, now.Location())

		// If birthday already passed this year, check next year
		if nextBday.Before(now) {
			nextBday = nextBday.AddDate(1, 0, 0)
		}

		// Check if within next N days
		daysUntil := int(nextBday.Sub(now).Hours() / 24)
		if daysUntil <= days {
			upcoming = append(upcoming, contact)
		}
	}

	return upcoming, nil
}

func (r *ContactRepo) GetNeedsAttention(userID uint, daysThreshold int) ([]entity.Contact, error) {
	var contacts []entity.Contact
	threshold := time.Now().AddDate(0, 0, -daysThreshold)

	// Get VIP contacts where last_contacted is older than threshold or empty
	err := r.DB.Where("user_id = ? AND vip = ? AND (last_contacted = '' OR last_contacted < ?)",
		userID, true, threshold.Format("2006-01-02")).
		Find(&contacts).Error

	return contacts, err
}

func (r *ContactRepo) GetRecentActivity(userID uint, limit int) ([]entity.Contact, error) {
	var contacts []entity.Contact
	err := r.DB.Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Find(&contacts).Error
	return contacts, err
}

// Event methods
func (r *ContactRepo) CreateEvent(event *entity.Event) error {
	return r.DB.Create(event).Error
}

func (r *ContactRepo) GetEventsByContact(contactID, userID uint) ([]entity.Event, error) {
	var events []entity.Event
	err := r.DB.Where("contact_id = ? AND user_id = ?", contactID, userID).
		Order("event_date ASC").
		Find(&events).Error
	return events, err
}

func (r *ContactRepo) GetEventByID(eventID, userID uint) (entity.Event, error) {
	var event entity.Event
	err := r.DB.Where("id = ? AND user_id = ?", eventID, userID).First(&event).Error
	return event, err
}

func (r *ContactRepo) DeleteEvent(eventID, userID uint) error {
	return r.DB.Where("id = ? AND user_id = ?", eventID, userID).Delete(&entity.Event{}).Error
}

func (r *ContactRepo) UpdateEventGoogleID(eventID, userID uint, googleEventID string) error {
	return r.DB.Model(&entity.Event{}).
		Where("id = ? AND user_id = ?", eventID, userID).
		Update("google_calendar_event_id", googleEventID).Error
}

func (r *ContactRepo) GetUpcomingEvents(userID uint, days int) ([]entity.Event, error) {
	var events []entity.Event
	err := r.DB.Where("user_id = ?", userID).
		Order("event_date ASC").
		Find(&events).Error
	return events, err
}
