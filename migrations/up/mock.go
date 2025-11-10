package main

import (
	"fmt"
	"github.com/La002/personal-crm/config"
	"github.com/La002/personal-crm/internal/entity"
	"github.com/La002/personal-crm/internal/repository"
	"github.com/La002/personal-crm/pkg/logger"
)

func main() {
	cfg := config.NewConfig()
	l := logger.New(cfg.Log.Level)
	repo := repository.NewContactRepo(cfg, l)

	contacts := []entity.Contact{
		{
			Name:         "Alice Johnson",
			Relationship: entity.Friend,
			Industry:     "Technology",
			Company:      "TechCorp",
			Birthday:     "1990-04-15",
			VIP:          true,
			DetailInfo: entity.DetailInfo{
				Location: "New York",
				Notes:    "Met at a tech conference.",
				FamilyDetails: entity.FamilyDetails{
					Spouse:   "John Johnson",
					Children: "Emily, Michael",
				},
				ContactInfo: entity.ContactInfo{
					PhoneNumber: "123-456-7890",
					Email:       "alice.johnson@example.com",
					LinkedIn:    "linkedin.com/in/alicejohnson",
					Instagram:   "instagram.com/alicejohnson",
				},
			},
			VIPInfo: entity.VIPInfo{
				LastMet:       "2023-08-15",
				LastContacted: "2023-09-05",
				LastUpdate:    "Gone to travel Japan",
			},
		},
		{
			Name:         "Bob Smith",
			Relationship: entity.Colleague,
			Industry:     "Finance",
			Company:      "FinCorp",
			Birthday:     "1985-09-20",
			VIP:          false,
		},
		{
			Name:         "Charlie Brown",
			Relationship: entity.Family,
			Industry:     "Healthcare",
			Company:      "HealthInc",
			Birthday:     "1992-12-05",
			VIP:          true,
			DetailInfo: entity.DetailInfo{
				Location: "Los Angeles",
				Notes:    "Close family member.",
				FamilyDetails: entity.FamilyDetails{
					Spouse:   "Sarah Brown",
					Children: "Liam, Sophia",
				},
				ContactInfo: entity.ContactInfo{
					PhoneNumber: "987-654-3210",
					Email:       "charlie.brown@example.com",
					LinkedIn:    "linkedin.com/in/charliebrown",
					Instagram:   "instagram.com/charliebrown",
				},
			},
			VIPInfo: entity.VIPInfo{
				LastMet:       "2023-08-15",
				LastContacted: "2023-09-05",
				LastUpdate:    "Changed Job",
			},
		},
		{
			Name:         "Diana Prince",
			Relationship: entity.Network,
			Industry:     "Education",
			Company:      "EduWorld",
			Birthday:     "1988-07-30",
			VIP:          false,
		},
	}

	for _, contact := range contacts {
		if err := repo.DB.Create(&contact).Error; err != nil {
			l.Error("Failed to insert contact: %v", err)
		}
	}

	fmt.Println("Mock data successfully✅ inserted.")
}
