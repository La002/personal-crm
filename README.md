# Personal CRM

A lightweight personal CRM system built with Go, featuring Google Calendar integration and comprehensive observability through Prometheus metrics.

## Features

- **Contact Management**: Store and manage personal contacts with rich metadata (relationships, industry, birthday, social links)
- **Google OAuth Authentication**: Secure login with Google accounts
- **Calendar Integration**: Sync birthdays and custom events to Google Calendar
- **Dashboard**: Quick overview of contacts and recent activities
- **Modern Frontend**: HTMX for dynamic interactions without JavaScript complexity + Tailwind CSS for responsive styling
- **Production-Ready Observability**:
  - Prometheus metrics for HTTP, database, and business metrics
  - Exposed via `/metrics` endpoint for scraping
  - Ready for integration with any Prometheus-compatible monitoring system
- **Database Migrations**: Version-controlled schema migrations using golang-migrate

## Tech Stack

- **Backend**: Go with Echo framework
- **Frontend**: HTMX + Tailwind CSS 
- **Database**: PostgreSQL 
- **Observability**: Prometheus metrics
- **Authentication**: Google OAuth 2.0
- **Deployment**: Render (app) + Neon (database)

## Quick Start

### Prerequisites
- Go 1.21+
- Docker (for local PostgreSQL)
- Google OAuth credentials

### Local Development

1. **Clone and install dependencies**
   ```bash
   git clone <repository-url>
   cd personal-crm
   go mod download
   ```

2. **Start PostgreSQL**
   ```bash
   docker run --name crm-postgres -e POSTGRES_PASSWORD=12345 -p 5432:5432 -d postgres:16
   ```

3. **Configure secrets**

   Create `config/config.secrets.yml`:
   ```yaml
   db:
     password: '12345'

   oauth:
     google_client_id: 'your-client-id'
     google_client_secret: 'your-client-secret'
     redirect_url: 'http://localhost:8080/auth/google/callback'

   jwt:
     secret_key: 'your-secret-key-min-32-chars'
   ```

4. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

5. **Access the app**
   - Application: http://localhost:8080
   - Metrics: http://localhost:8080/metrics

## Live Demo

The application is deployed and available to try out:

- **URL**: https://personal-crm-l8id.onrender.com
- **Metrics**: https://personal-crm-l8id.onrender.com/metrics

Deployed on Render (application) + Neon (PostgreSQL database) with zero-cost free tier infrastructure.

## Observability

The application exposes Prometheus metrics at `/metrics`:

- **HTTP Metrics**: Request counts, duration, in-flight requests
- **Database Metrics**: Connection pool stats (open, idle, in-use)
- **Business Metrics**: Total users, contacts, login attempts

Metrics can be scraped by any Prometheus-compatible monitoring system (Prometheus, Datadog, New Relic, etc.)

## Project Structure

```
.
├── cmd/server/          # Application entry point
├── config/              # Configuration management
├── internal/
│   ├── entity/         # Domain models
│   ├── middleware/     # HTTP middleware
│   ├── repository/     # Data access layer
│   ├── service/        # Business logic
│   ├── renderer/       # Template rendering
│   └── templates/      # HTML templates
├── migrations/         # Database migrations
├── pkg/
│   ├── logger/        # Logging utilities
│   ├── metrics/       # Prometheus metrics
│   └── postgres/      # Database connection
└── static/            # Static assets
```

## Database Migrations

Migrations run automatically on application startup using golang-migrate. Migration files are in `migrations/`:

- `000001_create_users_table.up.sql`
- `000002_create_contacts_table.up.sql`
- `000003_create_events_table.up.sql`

## Security

- OAuth credentials should be set via environment variables in production
- `config/config.yml` contains only dummy values (safe to commit)
- Real secrets go in `config/config.secrets.yml` (gitignored)
- JWT tokens with configurable expiry
- Rate limiting middleware (20 req/s)

