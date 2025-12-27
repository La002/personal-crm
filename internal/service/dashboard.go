package service

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/La002/personal-crm/internal/repository"
	"github.com/labstack/echo/v4"
)

type DashboardService struct {
	Repo *repository.ContactRepo
}

func NewDashboardService(repo *repository.ContactRepo) *DashboardService {
	return &DashboardService{
		Repo: repo,
	}
}

type DashboardData struct {
	UpcomingEvents []EventInfo
	NeedsAttention []AttentionInfo
	RecentActivity []ActivityInfo
}

type EventInfo struct {
	ID          uint
	ContactID   uint
	Name        string
	EventType   string // "Birthday" or "Custom"
	Title       string // Event title (for custom events)
	EventDate   string
	DaysUntil   int
	DisplayDate string
}

type AttentionInfo struct {
	ID            uint
	Name          string
	Company       string
	LastContacted string
	DaysSince     int
}

type ActivityInfo struct {
	ID        uint
	Name      string
	Company   string
	Action    string
	Timestamp string
}

func (s *DashboardService) GetDashboard(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	now := time.Now()
	allEvents := []EventInfo{}

	// Get upcoming birthdays (next 30 days)
	birthdayContacts, err := s.Repo.GetUpcomingBirthdays(userID, 30)
	if err != nil {
		c.Logger().Error("Failed to get birthdays: ", err)
	}

	for _, contact := range birthdayContacts {
		if contact.Birthday == "" {
			continue
		}

		bday, err := time.Parse("2006-01-02", contact.Birthday)
		if err != nil {
			continue
		}

		// Calculate next birthday
		nextBday := time.Date(now.Year(), bday.Month(), bday.Day(), 0, 0, 0, 0, now.Location())
		if nextBday.Before(now) {
			nextBday = nextBday.AddDate(1, 0, 0)
		}

		daysUntil := int(nextBday.Sub(now).Hours() / 24)

		allEvents = append(allEvents, EventInfo{
			ID:          contact.ID,
			ContactID:   contact.ID,
			Name:        contact.Name,
			EventType:   "Birthday",
			Title:       "",
			EventDate:   contact.Birthday,
			DaysUntil:   daysUntil,
			DisplayDate: nextBday.Format("Jan 2"),
		})
	}

	// Get custom events
	customEvents, err := s.Repo.GetUpcomingEvents(userID, 30)
	if err != nil {
		c.Logger().Error("Failed to get custom events: ", err)
	}

	for _, event := range customEvents {
		eventDate, err := time.Parse("2006-01-02", event.EventDate)
		if err != nil {
			continue
		}

		// Calculate days until event
		daysUntil := int(eventDate.Sub(now).Hours() / 24)

		// Only include events in the next 30 days
		if daysUntil < 0 || daysUntil > 30 {
			continue
		}

		// Get contact name
		contact, err := s.Repo.GetContact(fmt.Sprintf("%d", event.ContactID), userID)
		if err != nil {
			continue
		}

		allEvents = append(allEvents, EventInfo{
			ID:          event.ID,
			ContactID:   event.ContactID,
			Name:        contact.Name,
			EventType:   "Custom",
			Title:       event.Title,
			EventDate:   event.EventDate,
			DaysUntil:   daysUntil,
			DisplayDate: eventDate.Format("Jan 2"),
		})
	}

	// Get contacts needing attention (60+ days)
	attentionContacts, err := s.Repo.GetNeedsAttention(userID, 60)
	if err != nil {
		c.Logger().Error("Failed to get needs attention: ", err)
	}

	attention := []AttentionInfo{}
	for _, contact := range attentionContacts {
		daysSince := 0
		if contact.LastContacted != "" {
			lastContact, err := time.Parse("2006-01-02", contact.LastContacted)
			if err == nil {
				daysSince = int(now.Sub(lastContact).Hours() / 24)
			}
		}

		attention = append(attention, AttentionInfo{
			ID:            contact.ID,
			Name:          contact.Name,
			Company:       contact.Company,
			LastContacted: contact.LastContacted,
			DaysSince:     daysSince,
		})
	}

	// Get recent activity
	recentContacts, err := s.Repo.GetRecentActivity(userID, 10)
	if err != nil {
		c.Logger().Error("Failed to get recent activity: ", err)
	}

	activity := []ActivityInfo{}
	for _, contact := range recentContacts {
		action := "Updated"
		if contact.CreatedAt.Equal(contact.UpdatedAt) {
			action = "Added"
		}

		// Calculate relative time
		duration := now.Sub(contact.UpdatedAt)
		timeAgo := ""
		if duration.Hours() < 1 {
			timeAgo = "Just now"
		} else if duration.Hours() < 24 {
			hours := int(duration.Hours())
			if hours == 1 {
				timeAgo = "1 hour ago"
			} else {
				timeAgo = fmt.Sprintf("%d hours ago", hours)
			}
		} else {
			days := int(duration.Hours() / 24)
			if days == 1 {
				timeAgo = "1 day ago"
			} else {
				timeAgo = fmt.Sprintf("%d days ago", days)
			}
		}

		activity = append(activity, ActivityInfo{
			ID:        contact.ID,
			Name:      contact.Name,
			Company:   contact.Company,
			Action:    action,
			Timestamp: timeAgo,
		})
	}

	// Sort all events by days until
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].DaysUntil < allEvents[j].DaysUntil
	})

	data := DashboardData{
		UpcomingEvents: allEvents,
		NeedsAttention: attention,
		RecentActivity: activity,
	}

	return c.Render(http.StatusOK, "dashboard", data)
}
