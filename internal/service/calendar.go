package service

import (
	"context"
	"fmt"

	"github.com/La002/personal-crm/internal/entity"
	"github.com/La002/personal-crm/internal/repository"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarService struct {
	UserRepo     repository.UserDao
	ContactRepo  repository.ContactDao
	OAuth2Config *oauth2.Config
}

func NewCalendarService(
	userRepo repository.UserDao,
	contactRepo repository.ContactDao,
	clientID string,
	clientSecret string,
	redirectURL string,
) *CalendarService {
	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/calendar.events",
		},
		Endpoint: google.Endpoint,
	}
	return &CalendarService{
		UserRepo:     userRepo,
		ContactRepo:  contactRepo,
		OAuth2Config: oauth2Config,
	}
}

func (s *CalendarService) CreateBirthdayReminder(userID uint, contactID string) error {
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user")
	}
	contact, err := s.ContactRepo.GetContact(contactID, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch contact")
	}

	client, err := s.getCalendarClient(&user)
	if err != nil {
		return fmt.Errorf("failed to fetch client")
	}

	event := &calendar.Event{
		Summary: fmt.Sprintf("%s's Birthday", contact.Name),
		Start: &calendar.EventDateTime{
			Date: contact.Birthday,
		},
		End: &calendar.EventDateTime{
			Date: contact.Birthday,
		},
		Recurrence: []string{"RRULE:FREQ=YEARLY"},
	}

	createdEvent, err := client.Events.Insert("primary", event).Do()
	if err != nil {
		return fmt.Errorf("failed to create calendar event: %w", err)
	}

	return s.ContactRepo.UpdateCalendarSync(contactID, userID, createdEvent.Id, true)
}
func (s *CalendarService) DeleteBirthdayReminder(userID uint, contactID string) error {
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user")
	}

	contact, err := s.ContactRepo.GetContact(contactID, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch contact")
	}
	client, err := s.getCalendarClient(&user)
	if err != nil {
		return err
	}

	// Delete from Google Calendar
	err = client.Events.Delete("primary", contact.GoogleCalendarEventID).Do()
	if err != nil {
		return err
	}

	// Update database to clear sync status
	return s.ContactRepo.UpdateCalendarSync(contactID, userID, "", false)
}
func (s *CalendarService) GetEventStatus(userID uint, contactID string) (bool, error) {
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch user")
	}
	contact, err := s.ContactRepo.GetContact(contactID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch contact")
	}

	// If no calendar event ID is stored, the event doesn't exist
	if contact.GoogleCalendarEventID == "" {
		return false, nil
	}

	client, err := s.getCalendarClient(&user)
	if err != nil {
		return false, err
	}

	// Try to fetch the event from Google Calendar
	event, err := client.Events.Get("primary", contact.GoogleCalendarEventID).Do()
	if err != nil {
		// Event not found or deleted, return false
		return false, nil
	}

	// Check if event exists and is not cancelled
	if event != nil && event.Status != "cancelled" {
		return true, nil
	}

	return false, nil
}

// Custom event methods
func (s *CalendarService) CreateCustomEvent(userID, contactID uint, title, eventDate, recurrence string) error {
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user")
	}

	// Get contact to prefix name
	contact, err := s.ContactRepo.GetContact(fmt.Sprintf("%d", contactID), userID)
	if err != nil {
		return fmt.Errorf("failed to fetch contact")
	}

	client, err := s.getCalendarClient(&user)
	if err != nil {
		return fmt.Errorf("failed to get calendar client")
	}

	// Prefix contact's name to event title
	fullTitle := fmt.Sprintf("%s - %s", contact.Name, title)

	// Create event in Google Calendar
	calEvent := &calendar.Event{
		Summary: fullTitle,
		Start: &calendar.EventDateTime{
			Date: eventDate,
		},
		End: &calendar.EventDateTime{
			Date: eventDate,
		},
	}

	// Add recurrence if specified
	if recurrence == "monthly" {
		calEvent.Recurrence = []string{"RRULE:FREQ=MONTHLY"}
	} else if recurrence == "yearly" {
		calEvent.Recurrence = []string{"RRULE:FREQ=YEARLY"}
	}

	createdEvent, err := client.Events.Insert("primary", calEvent).Do()
	if err != nil {
		return fmt.Errorf("failed to create calendar event: %w", err)
	}

	// Save event to database
	event := &entity.Event{
		UserID:                userID,
		ContactID:             contactID,
		Title:                 title, // Store original title without prefix
		EventDate:             eventDate,
		Recurrence:            recurrence,
		GoogleCalendarEventID: createdEvent.Id,
	}

	return s.ContactRepo.CreateEvent(event)
}

func (s *CalendarService) DeleteCustomEvent(userID, eventID uint) error {
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user")
	}

	event, err := s.ContactRepo.GetEventByID(eventID, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch event")
	}

	if event.GoogleCalendarEventID == "" {
		// Just delete from database
		return s.ContactRepo.DeleteEvent(eventID, userID)
	}

	client, err := s.getCalendarClient(&user)
	if err != nil {
		return fmt.Errorf("failed to get calendar client")
	}

	// Delete from Google Calendar
	err = client.Events.Delete("primary", event.GoogleCalendarEventID).Do()
	if err != nil {
		return fmt.Errorf("failed to delete calendar event: %w", err)
	}

	// Delete from database
	return s.ContactRepo.DeleteEvent(eventID, userID)
}

func (s *CalendarService) getCalendarClient(user *entity.User) (*calendar.Service, error) {
	ctx := context.Background()

	token := oauth2.Token{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		Expiry:       user.TokenExpiry,
	}

	// Refresh token if expired
	if !token.Valid() {
		if token.RefreshToken == "" {
			return nil, fmt.Errorf("access token expired and no refresh token available - user must re-authenticate")
		}

		if err := s.refreshAccessToken(user); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		// Update token with refreshed values
		token.AccessToken = user.AccessToken
		token.RefreshToken = user.RefreshToken
		token.Expiry = user.TokenExpiry
	}

	client := s.OAuth2Config.Client(ctx, &token)

	service, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	return service, nil
}

func (s *CalendarService) refreshAccessToken(user *entity.User) error {
	ctx := context.Background()

	token := &oauth2.Token{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		Expiry:       user.TokenExpiry,
	}

	// Use the OAuth2 config to get a new token
	tokenSource := s.OAuth2Config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return fmt.Errorf("failed to refresh access token: %w", err)
	}

	// Update user with new token details
	user.AccessToken = newToken.AccessToken
	user.RefreshToken = newToken.RefreshToken
	user.TokenExpiry = newToken.Expiry

	return s.UserRepo.UpdateUser(user)
}
