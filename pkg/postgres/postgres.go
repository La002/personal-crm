package postgres

import (
	"fmt"
	"github.com/La002/personal-crm/config"
	"github.com/La002/personal-crm/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(config *config.Configuration, l *logger.Logger) *gorm.DB {
	connectStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.Name)
	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{})
	if err != nil {
		l.Fatal(err)
	}
	return db
}
