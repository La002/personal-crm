package main

import (
	"time"

	"github.com/La002/personal-crm/config"
	"github.com/La002/personal-crm/internal/middleware"
	"github.com/La002/personal-crm/internal/renderer"
	"github.com/La002/personal-crm/internal/service"
	"github.com/La002/personal-crm/migrations"
	"github.com/La002/personal-crm/pkg/logger"
	"github.com/La002/personal-crm/pkg/metrics"
	"github.com/La002/personal-crm/pkg/postgres"
	"github.com/La002/personal-crm/pkg/repository"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

func main() {
	// Initialize configuration and logger
	cfg := config.NewConfig()
	l := logger.New(cfg.Log.Level)
	gormDB := postgres.ConnectDB(cfg, l)

	sqlDB, err := gormDB.DB()
	if err != nil {
		l.Fatal("Failed to get underlying SQL DB: %w", err)
	}

	if err = sqlDB.Ping(); err != nil {
		l.Fatal("Database ping failed: %w", err)
	}

	l.Debug("Running database migrations...")
	if err = migrations.RunMigrations(sqlDB); err != nil {
		l.Fatal("Failed to run migrations:", err)
	}
	l.Debug("Migrations completed successfully")

	// Start background metrics updater
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			metrics.UpdateDBMetrics(gormDB)
			metrics.UpdateBusinessMetrics(gormDB)
		}
	}()

	// Initialize repositories
	contactRepo := repository.NewContactRepo(cfg, l)
	userRepo := repository.NewUserRepo(cfg, l)

	// Initialize services
	contactService := service.NewContactService(contactRepo)
	dashboardService := service.NewDashboardService(contactRepo)

	// Debug: Log OAuth configuration
	//clientSecretPreview := "EMPTY"
	//if len(cfg.OAuth.GoogleClientSecret) > 0 {
	//	clientSecretPreview = cfg.OAuth.GoogleClientSecret[:10] + "..."
	//}
	//l.Info("OAuth Config",
	//	"ClientID", cfg.OAuth.GoogleClientID,
	//	"ClientIDEmpty", cfg.OAuth.GoogleClientID == "",
	//	"ClientSecret", clientSecretPreview,
	//	"RedirectURL", cfg.OAuth.RedirectURL)

	authService := service.NewAuthService(
		userRepo,
		cfg.OAuth.GoogleClientID,
		cfg.OAuth.GoogleClientSecret,
		cfg.OAuth.RedirectURL,
		cfg.JWT.SecretKey,
		cfg.JWT.ExpiryHours,
	)

	calendarService := service.NewCalendarService(
		userRepo,
		contactRepo,
		cfg.OAuth.GoogleClientID,
		cfg.OAuth.GoogleClientSecret,
		cfg.OAuth.RedirectURL)

	authHandler := service.NewAuthHandler(authService)
	calendarHandler := service.NewCalendarHandler(calendarService)
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		l.Error("HTTP Error occurred", "error", err.Error())
		c.JSON(500, map[string]string{
			"error": err.Error(),
		})
	}

	// Add logger middleware
	e.Use(echomiddleware.LoggerWithConfig(echomiddleware.LoggerConfig{
		Format: "time=${time_rfc3339} method=${method} uri=${uri} status=${status} error=${error}\n",
	}))

	// Add Prometheus metrics middleware
	e.Use(metrics.PrometheusMiddleware())

	// Little bit of middlewares for housekeeping
	e.Pre(echomiddleware.RemoveTrailingSlash())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RateLimiter(echomiddleware.NewRateLimiterMemoryStore(
		rate.Limit(20),
	)))
	e.Static("/static", "static")
	renderer.NewTemplateRenderer(e, "internal/templates/*.html")

	// Public routes (no authentication required)
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "home", nil)
	})

	// Prometheus metrics endpoint
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/auth/google/login", authHandler.GoogleLogin)
	e.GET("/auth/google/callback", authHandler.GoogleCallback)
	e.GET("/login", func(c echo.Context) error {
		return c.Render(200, "login", nil)
	})

	// Protected routes (authentication required)
	protected := e.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))

	// Dashboard
	protected.GET("/dashboard", dashboardService.GetDashboard)

	// Contacts
	protected.GET("/contacts", contactService.GetAllContacts)
	protected.GET("/contacts/search", contactService.SearchContacts)
	protected.POST("/contacts/new", contactService.AddContact)
	protected.DELETE("/contacts/:id", contactService.DeleteContact)
	protected.GET("/contacts/:id", contactService.GetContact)
	protected.GET("/contacts/:id/edit", contactService.EditContact)
	protected.PUT("/contacts/:id", contactService.UpdateContact)
	protected.POST("/auth/logout", authHandler.Logout)

	// Calendar sync endpoints
	protected.POST("/contacts/:id/calendar/sync", calendarHandler.SyncBirthdayToCalendar)
	protected.DELETE("/contacts/:id/calendar/sync", calendarHandler.DeleteCalendarSync)
	protected.GET("/contacts/:id/calendar/status", calendarHandler.GetCalendarSyncStatus)

	// Custom events endpoints
	protected.POST("/contacts/:id/events", calendarHandler.CreateCustomEvent)
	protected.DELETE("/contacts/:id/events/:eventId", calendarHandler.DeleteCustomEvent)

	// Start the HTTP server
	e.Logger.Fatal(e.Start(":8080"))
}
