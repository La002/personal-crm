package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type CalendarHandler struct {
	CalendarService *CalendarService
}

func NewCalendarHandler(calendarService *CalendarService) *CalendarHandler {
	return &CalendarHandler{
		CalendarService: calendarService,
	}
}

func (h *CalendarHandler) SyncBirthdayToCalendar(c echo.Context) error {
	contactID := c.Param("id")
	userID := c.Get("user_id").(uint)

	err := h.CalendarService.CreateBirthdayReminder(userID, contactID)
	if err != nil {
		c.Logger().Error("Failed to create calendar reminder: ", err)

		// Check if it's an authentication error
		errMsg := err.Error()
		if strings.Contains(errMsg, "must re-authenticate") || strings.Contains(errMsg, "refresh token") {
			return c.String(401, "Authentication expired. Please logout and login again.")
		}

		return c.String(500, "Failed to sync calendar. Please try again.")
	}

	// Fetch updated contact
	contact, err := h.CalendarService.ContactRepo.GetContact(contactID, userID)
	if err != nil {
		return c.String(500, "Failed to fetch updated contact")
	}

	// Return single row HTML
	contactData := getContactMapShort(contact)
	return c.Render(http.StatusOK, "contact-row", contactData)
}

func (h *CalendarHandler) DeleteCalendarSync(c echo.Context) error {
	contactID := c.Param("id")
	userID := c.Get("user_id").(uint)

	err := h.CalendarService.DeleteBirthdayReminder(userID, contactID)
	if err != nil {
		c.Logger().Error("Failed to delete calendar reminder: ", err)

		// Check if it's an authentication error
		errMsg := err.Error()
		if strings.Contains(errMsg, "must re-authenticate") || strings.Contains(errMsg, "refresh token") {
			return c.String(401, "Authentication expired. Please logout and login again.")
		}

		return c.String(500, "Failed to delete calendar sync. Please try again.")
	}

	// Fetch updated contact
	contact, err := h.CalendarService.ContactRepo.GetContact(contactID, userID)
	if err != nil {
		return c.String(500, "Failed to fetch updated contact")
	}

	// Return single row HTML
	contactData := getContactMapShort(contact)
	return c.Render(http.StatusOK, "contact-row", contactData)
}

func (h *CalendarHandler) GetCalendarSyncStatus(c echo.Context) error {
	contactID := c.Param("id")
	userID := c.Get("user_id").(uint)

	synced, err := h.CalendarService.GetEventStatus(userID, contactID)
	if err != nil {
		c.Logger().Error("Failed to get calendar sync status: ", err)
		return c.String(500, "Failed to get sync status")
	}

	return c.JSON(200, map[string]bool{"synced": synced})
}

func (h *CalendarHandler) CreateCustomEvent(c echo.Context) error {
	contactID := c.Param("id")
	userID := c.Get("user_id").(uint)

	title := c.FormValue("title")
	eventDate := c.FormValue("event_date")
	recurrence := c.FormValue("recurrence")

	if title == "" || eventDate == "" {
		return c.String(400, "Title and event date are required")
	}

	// Default to "none" if not specified
	if recurrence == "" {
		recurrence = "none"
	}

	// Parse contactID to uint
	var cID uint
	_, err := fmt.Sscan(contactID, &cID)
	if err != nil {
		return c.String(400, "Invalid contact ID")
	}

	err = h.CalendarService.CreateCustomEvent(userID, cID, title, eventDate, recurrence)
	if err != nil {
		c.Logger().Error("Failed to create custom event: ", err)

		// Check if it's an authentication error
		errMsg := err.Error()
		if strings.Contains(errMsg, "must re-authenticate") || strings.Contains(errMsg, "refresh token") {
			return c.String(401, "Authentication expired. Please logout and login again.")
		}

		return c.String(500, "Failed to create event. Please try again.")
	}

	// Fetch the newly created event to render
	events, err := h.CalendarService.ContactRepo.GetEventsByContact(cID, userID)
	if err != nil || len(events) == 0 {
		return c.String(500, "Failed to fetch created event")
	}

	// Get the last event (the one we just created)
	newEvent := events[len(events)-1]

	// Return just the event row HTML for HTMX to append
	return c.Render(http.StatusOK, "event-row", newEvent)
}

func (h *CalendarHandler) DeleteCustomEvent(c echo.Context) error {
	eventID := c.Param("eventId")
	userID := c.Get("user_id").(uint)

	// Parse eventID to uint
	var eID uint
	_, err := fmt.Sscan(eventID, &eID)
	if err != nil {
		return c.String(400, "Invalid event ID")
	}

	err = h.CalendarService.DeleteCustomEvent(userID, eID)
	if err != nil {
		c.Logger().Error("Failed to delete custom event: ", err)
		return c.String(500, "Failed to delete event")
	}

	// Return empty response for HTMX to swap out the element
	return c.NoContent(http.StatusOK)
}
