package main

import (
	"flag"
	"log"
	"os"

	"github.com/La002/personal-crm/config"
	"github.com/La002/personal-crm/internal/repository"
	"github.com/La002/personal-crm/migrations"
	"github.com/La002/personal-crm/pkg/logger"
)

func main() {
	// Parse number of steps to rollback
	steps := flag.Int("steps", 1, "Number of migrations to rollback")
	flag.Parse()

	cfg := config.NewConfig()
	l := logger.New(cfg.Log.Level)
	repo := repository.NewContactRepo(cfg, l)

	// Extract *sql.DB from GORM
	sqlDB, err := repo.DB.DB()
	if err != nil {
		log.Fatal("Failed to get DB:", err)
	}

	l.Info("Rolling back %d migration(s)...", *steps)

	err = migrations.RollbackMigrations(sqlDB, *steps)
	if err != nil {
		l.Error("Rollback failed: %s", err)
		os.Exit(1)
	}

	l.Info("Rollback completed successfully!")
}
