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
	repo.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	repo.DB.Exec("CREATE TYPE relation AS ENUM ('Friend', 'Family', 'Colleague', 'School', 'Network', 'Services')")
	err := repo.DB.AutoMigrate(
		&entity.Contact{},
		&entity.DetailChanges{})
	if err != nil {
		l.Fatal("Auto migration failed")
	}

	fmt.Println("👍 Migration complete")
}
