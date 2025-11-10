package main

import (
	"github.com/La002/personal-crm/config"
	"github.com/La002/personal-crm/internal/entity"
	"github.com/La002/personal-crm/internal/repository"
	"github.com/La002/personal-crm/pkg/logger"
)

func main() {
	cfg := config.NewConfig()
	l := logger.New(cfg.Log.Level)
	repo := repository.NewContactRepo(cfg, l)
	err := repo.DB.Migrator().DropTable(&entity.Contact{}, &entity.DetailChanges{})
	if err != nil {
		l.Error("error happened: %s", err)
		return
	}
}
