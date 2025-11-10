package main

import (
	"github.com/La002/personal-crm/config"
	"github.com/La002/personal-crm/internal/renderer"
	"github.com/La002/personal-crm/internal/repository"
	"github.com/La002/personal-crm/internal/service"
	"github.com/La002/personal-crm/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/time/rate"
)

func main() {
	// Initialize configuration and logger
	cfg := config.NewConfig()
	l := logger.New(cfg.Log.Level)

	// Initialize repository
	contactRepo := repository.NewContactRepo(cfg, l)
	service := service.NewContactService(contactRepo)
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		l.Error("HTTP Error occurred", "error", err.Error())
		c.JSON(500, map[string]string{
			"error": err.Error(),
		})
	}

	// Add logger middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339} method=${method} uri=${uri} status=${status} error=${error}\n",
	}))

	// Little bit of middlewares for housekeeping
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(20),
	)))
	e.Static("/static", "static")
	renderer.NewTemplateRenderer(e, "internal/templates/*.html")

	e.GET("/contacts", service.GetAllContacts)
	e.POST("/contacts/new", service.AddContact)
	e.DELETE("/contacts/:id", service.DeleteContact)
	e.GET("contacts/:id", service.GetContact)
	e.GET("/contacts/:id/edit", service.EditContact)
	e.PUT("/contacts/:id", service.UpdateContact)

	// Start the HTTP server
	e.Logger.Fatal(e.Start(":8080"))
}
