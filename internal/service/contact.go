package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/La002/personal-crm/internal/entity"
	"github.com/La002/personal-crm/internal/repository"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ContactService struct {
	Repo repository.ContactDao
}

func NewContactService(repo repository.ContactDao) *ContactService {
	return &ContactService{
		Repo: repo,
	}
}

// GetAllContacts retrieves all contacts and renders them as HTML for HTMX
func (s *ContactService) GetAllContacts(c echo.Context) error {
	var res []map[string]interface{}
	userID := c.Get("user_id").(uint)

	contacts, err := s.Repo.GetAllContacts(userID)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		mp := getContactMapShort(contact)
		res = append(res, mp)
	}

	// Get user email from context (set by middleware)
	email := ""
	if emailVal := c.Get("user_email"); emailVal != nil {
		email, _ = emailVal.(string)
	}

	data := map[string]interface{}{
		"Contacts": res,
		"Email":    email,
	}
	return c.Render(http.StatusOK, "contacts", data)
}

func (s *ContactService) SearchContacts(c echo.Context) error {
	filter := c.QueryParam("filter")
	userID := c.Get("user_id").(uint)

	var contacts []entity.Contact
	var err error
	switch filter {
	case "vip":
		contacts, err = s.Repo.FilterContactsByVipStatus(true, userID)
	case "non-vip":
		contacts, err = s.Repo.FilterContactsByVipStatus(false, userID)
	case "Friend", "Family", "Colleague", "School", "Network", "Services":
		contacts, err = s.Repo.SearchContactsByRelationship(filter, userID)
	default:
		contacts, err = s.Repo.GetAllContacts(userID)
	}

	if err != nil {
		return err
	}
	var res []map[string]interface{}
	for _, contact := range contacts {
		mp := getContactMapShort(contact)
		res = append(res, mp)
	}

	return c.Render(http.StatusOK, "table-content", res)
}

func (s *ContactService) AddContact(c echo.Context) error {
	vipValue := c.FormValue("vip") == "on"
	userID := c.Get("user_id").(uint)

	newContact := &entity.Contact{
		Name:         c.FormValue("name"),
		Relationship: entity.Relation(c.FormValue("relationship")),
		Industry:     c.FormValue("industry"),
		Company:      c.FormValue("company"),
		Birthday:     c.FormValue("birthday"),
		Vip:          vipValue,
		UserID:       userID,
	}
	if vipValue {
		newContact.LastContacted = c.FormValue("last_contacted")
		newContact.LastMet = c.FormValue("last_met")
		newContact.LastUpdate = c.FormValue("last_update")
	}

	if err := s.Repo.CreateContact(newContact); err != nil {
		return err
	}

	var res []map[string]interface{}

	contacts, err := s.Repo.GetAllContacts(userID)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		mp := getContactMapShort(contact)
		res = append(res, mp)
	}

	return c.Render(http.StatusOK, "table-content", res)
}

func (s *ContactService) DeleteContact(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("user_id").(uint)

	if err := s.Repo.DeleteContact(id, userID); err != nil {
		return err
	}

	contacts, err := s.Repo.GetAllContacts(userID)
	if err != nil {
		return err
	}

	var res []map[string]interface{}
	for _, contact := range contacts {
		mp := getContactMapShort(contact)
		res = append(res, mp)
	}

	return c.Render(http.StatusOK, "table-content", res)
}

func (s *ContactService) GetContact(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("user_id").(uint)

	contact, err := s.Repo.GetContact(id, userID)
	if err != nil {
		return err
	}

	// Fetch events for this contact
	events, err := s.Repo.GetEventsByContact(contact.ID, userID)
	if err != nil {
		c.Logger().Error("Failed to fetch events: ", err)
		events = []entity.Event{} // Empty events on error
	}

	res := getContactMapLong(contact)
	res["Events"] = events
	return c.Render(http.StatusOK, "detail", res)
}

func (s *ContactService) EditContact(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("user_id").(uint)
	contact, err := s.Repo.GetContact(id, userID)
	if err != nil {
		return err
	}

	res := getContactMapLong(contact)
	return c.Render(http.StatusOK, "edit", res)
}

// UpdateContact handles the form submission
func (s *ContactService) UpdateContact(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("user_id").(uint)
	vipValue := c.FormValue("vip") == "on"
	c.Logger().Info("update id : ", id)

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}
	contact := &entity.Contact{
		Model: gorm.Model{
			ID: uint(idUint),
		},
		UserID:       userID,
		Name:         c.FormValue("name"),
		Relationship: entity.Relation(c.FormValue("relationship")),
		Industry:     c.FormValue("industry"),
		Company:      c.FormValue("company"),
		Birthday:     c.FormValue("birthday"),
		Vip:          vipValue,
		DetailInfo: entity.DetailInfo{
			FamilyDetails: entity.FamilyDetails{
				Spouse:   c.FormValue("spouse"),
				Children: c.FormValue("children"),
			},
			ContactInfo: entity.ContactInfo{
				Location:    c.FormValue("location"),
				PhoneNumber: c.FormValue("phone"),
				Email:       c.FormValue("email"),
				LinkedIn:    c.FormValue("linkedin"),
				Instagram:   c.FormValue("instagram"),
				X:           c.FormValue("x"),
			},
			Notes: c.FormValue("notes"),
		},
	}

	if contact.Vip {
		contact.VipInfo = entity.VipInfo{
			LastContacted: c.FormValue("last_contacted"),
			LastMet:       c.FormValue("last_met"),
			LastUpdate:    c.FormValue("last_update"),
		}
	}

	if err := s.Repo.SaveContact(contact); err != nil {
		return err
	}

	res := getContactMapLong(*contact)

	return c.Render(http.StatusOK, "blocks", res)
}

func getContactMapShort(contact entity.Contact) map[string]interface{} {
	return map[string]interface{}{
		"Id":                    contact.ID,
		"Name":                  contact.Name,
		"Relationship":          contact.Relationship,
		"Industry":              contact.Industry,
		"Company":               contact.Company,
		"Birthday":              contact.Birthday,
		"Vip":                   contact.Vip,
		"LastContacted":         contact.LastContacted,
		"LastMet":               contact.LastMet,
		"LastUpdate":            contact.LastUpdate,
		"CalendarSyncEnabled":   contact.CalendarSyncEnabled,
		"GoogleCalendarEventID": contact.GoogleCalendarEventID,
	}
}

func getContactMapLong(contact entity.Contact) map[string]interface{} {
	res := map[string]interface{}{
		"Id":                    contact.ID,
		"Name":                  contact.Name,
		"Relationship":          contact.Relationship,
		"Industry":              contact.Industry,
		"Company":               contact.Company,
		"Birthday":              contact.Birthday,
		"Vip":                   contact.Vip,
		"Spouse":                contact.Spouse,
		"Children":              contact.Children,
		"Phone":                 contact.PhoneNumber,
		"Email":                 contact.Email,
		"LinkedIn":              contact.LinkedIn,
		"Instagram":             contact.Instagram,
		"X":                     contact.X,
		"Location":              contact.Location,
		"Notes":                 contact.Notes,
		"LastContacted":         contact.LastContacted,
		"LastMet":               contact.LastMet,
		"LastUpdate":            contact.LastUpdate,
		"CalendarSyncEnabled":   contact.CalendarSyncEnabled,
		"GoogleCalendarEventID": contact.GoogleCalendarEventID,
	}
	return res
}
