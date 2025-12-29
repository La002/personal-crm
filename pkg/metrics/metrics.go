package metrics

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"
)

// HTTP Metrics
var (
	// HTTPRequestsTotal tracks total HTTP requests by method, path, and status
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestDuration tracks request latency in seconds
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets, // Default buckets: 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
		},
		[]string{"method", "path"},
	)

	// HTTPRequestsInFlight tracks current number of requests being processed
	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)
)

// Database Metrics
var (
	// DBConnectionsOpen tracks open database connections
	DBConnectionsOpen = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_open",
			Help: "Number of open database connections",
		},
	)

	// DBConnectionsIdle tracks idle database connections
	DBConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	// DBConnectionsInUse tracks connections currently in use
	DBConnectionsInUse = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_in_use",
			Help: "Number of database connections currently in use",
		},
	)

	// DBQueryDuration tracks database query latency
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query latency in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"operation"},
	)
)

// Business Metrics
var (
	// ContactsTotal tracks total number of contacts
	ContactsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "contacts_total",
			Help: "Total number of contacts in the system",
		},
	)

	// ContactsCreated tracks contacts created (counter)
	ContactsCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "contacts_created_total",
			Help: "Total number of contacts created",
		},
	)

	// UsersTotal tracks total number of users
	UsersTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "users_total",
			Help: "Total number of registered users",
		},
	)

	// UserLoginsTotal tracks user login attempts
	UserLoginsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_logins_total",
			Help: "Total number of user login attempts",
		},
		[]string{"status"}, // success or failure
	)
)

// PrometheusMiddleware records HTTP metrics for each request
func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Increment in-flight requests
			HTTPRequestsInFlight.Inc()
			defer HTTPRequestsInFlight.Dec()

			// Process request
			err := next(c)

			// Record metrics
			duration := time.Since(start).Seconds()
			status := strconv.Itoa(c.Response().Status)
			method := c.Request().Method
			path := c.Path()

			// Don't track /metrics endpoint itself to avoid noise
			if path != "/metrics" {
				HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
				HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
			}

			return err
		}
	}
}

// UpdateDBMetrics updates database connection pool metrics
func UpdateDBMetrics(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	stats := sqlDB.Stats()
	DBConnectionsOpen.Set(float64(stats.OpenConnections))
	DBConnectionsIdle.Set(float64(stats.Idle))
	DBConnectionsInUse.Set(float64(stats.InUse))
}

// RecordLogin records a user login attempt
func RecordLogin(success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	UserLoginsTotal.WithLabelValues(status).Inc()
}

// RecordContactCreated increments the contact creation counter
func RecordContactCreated() {
	ContactsCreated.Inc()
}

// UpdateBusinessMetrics updates gauge metrics for users and contacts
func UpdateBusinessMetrics(db *gorm.DB) {
	var userCount, contactCount int64

	db.Model(&struct{ ID uint }{}).Table("users").Where("deleted_at IS NULL").Count(&userCount)
	db.Model(&struct{ ID uint }{}).Table("contacts").Where("deleted_at IS NULL").Count(&contactCount)

	UsersTotal.Set(float64(userCount))
	ContactsTotal.Set(float64(contactCount))
}
