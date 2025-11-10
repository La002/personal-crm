package service

import (
	"fmt"
	"github.com/La002/personal-crm/internal/entity"
	"github.com/La002/personal-crm/internal/repository"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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

	contacts, err := s.Repo.GetAllContacts()
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		mp := map[string]interface{}{
			"Id":            contact.ID,
			"Name":          contact.Name,
			"Relationship":  contact.Relationship,
			"Industry":      contact.Industry,
			"Company":       contact.Company,
			"Birthday":      contact.Birthday,
			"VIP":           contact.VIP,
			"LastContacted": contact.LastContacted,
			"LastMet":       contact.LastMet,
			"LastUpdate":    contact.LastUpdate,
		}
		res = append(res, mp)
	}
	return c.Render(http.StatusOK, "contacts", res)
}

func (s *ContactService) AddContact(c echo.Context) error {
	vipValue := c.FormValue("vip") == "on"
	newContact := &entity.Contact{
		Name:         c.FormValue("name"),
		Relationship: entity.Relation(c.FormValue("relationship")),
		Industry:     c.FormValue("industry"),
		Company:      c.FormValue("company"),
		Birthday:     c.FormValue("birthday"),
		VIP:          vipValue,
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

	contacts, err := s.Repo.GetAllContacts()
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		mp := map[string]interface{}{
			"Id":            contact.ID,
			"Name":          contact.Name,
			"Relationship":  contact.Relationship,
			"Industry":      contact.Industry,
			"Company":       contact.Company,
			"Birthday":      contact.Birthday,
			"VIP":           contact.VIP,
			"LastContacted": contact.LastContacted,
			"LastMet":       contact.LastMet,
			"LastUpdate":    contact.LastUpdate,
		}
		res = append(res, mp)
	}

	return c.Render(http.StatusOK, "table-content", res)
}

func (s *ContactService) DeleteContact(c echo.Context) error {
	id := c.Param("id")

	if err := s.Repo.DeleteContact(id); err != nil {
		return err
	}

	contacts, err := s.Repo.GetAllContacts()
	if err != nil {
		return err
	}

	var res []map[string]interface{}
	for _, contact := range contacts {
		mp := map[string]interface{}{
			"Id":            contact.ID,
			"Name":          contact.Name,
			"Relationship":  contact.Relationship,
			"Industry":      contact.Industry,
			"Company":       contact.Company,
			"Birthday":      contact.Birthday,
			"VIP":           contact.VIP,
			"LastContacted": contact.LastContacted,
			"LastMet":       contact.LastMet,
			"LastUpdate":    contact.LastUpdate,
		}
		res = append(res, mp)
	}

	return c.Render(http.StatusOK, "table-content", res)
}

func (s *ContactService) GetContact(c echo.Context) error {
	id := c.Param("id")

	contact, err := s.Repo.GetContact(id)
	if err != nil {
		return err
	}

	res := getContactMapLong(contact)
	return c.Render(http.StatusOK, "detail", res)
}

func (s *ContactService) EditContact(c echo.Context) error {
	id := c.Param("id")
	contact, err := s.Repo.GetContact(id)
	if err != nil {
		return err
	}

	res := getContactMapLong(contact)
	return c.Render(http.StatusOK, "edit", res)
}

// UpdateContact handles the form submission
func (s *ContactService) UpdateContact(c echo.Context) error {
	id := c.Param("id")
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
		Name:         c.FormValue("name"),
		Relationship: entity.Relation(c.FormValue("relationship")),
		Industry:     c.FormValue("industry"),
		Company:      c.FormValue("company"),
		Birthday:     c.FormValue("birthday"),
		VIP:          vipValue,
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

	if contact.VIP {
		contact.VIPInfo = entity.VIPInfo{
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

func getContactMapLong(contact entity.Contact) map[string]interface{} {
	res := map[string]interface{}{
		"Id":            contact.ID,
		"Name":          contact.Name,
		"Relationship":  contact.Relationship,
		"Industry":      contact.Industry,
		"Company":       contact.Company,
		"Birthday":      contact.Birthday,
		"VIP":           contact.VIP,
		"Spouse":        contact.Spouse,
		"Children":      contact.Children,
		"Phone":         contact.PhoneNumber,
		"Email":         contact.Email,
		"LinkedIn":      contact.LinkedIn,
		"Instagram":     contact.Instagram,
		"X":             contact.X,
		"Location":      contact.Location,
		"Notes":         contact.Notes,
		"LastContacted": contact.LastContacted,
		"LastMet":       contact.LastMet,
		"LastUpdate":    contact.LastUpdate,
	}
	return res
}
