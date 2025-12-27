package middleware

import (
	"net/http"

	jwtutil "github.com/La002/personal-crm/pkg/jwt"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware validates JWT tokens and extracts user information
func AuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Read JWT from cookie
			cookie, err := c.Cookie("auth_token")
			if err != nil {
				// No cookie found, redirect to login
				return c.Redirect(http.StatusFound, "/login")
			}

			// Validate JWT token
			claims, err := jwtutil.ValidateToken(cookie.Value, jwtSecret)
			if err != nil {
				// Invalid token, redirect to login
				return c.Redirect(http.StatusFound, "/login")
			}

			// Extract user ID from claims
			userID, err := jwtutil.ExtractUserID(claims)
			if err != nil {
				// Could not extract user ID, redirect to login
				return c.Redirect(http.StatusFound, "/login")
			}

			// Extract email from claims
			email, _ := (*claims)["email"].(string)

			// Store user ID and email in context for handlers to use
			c.Set("user_id", userID)
			c.Set("user_email", email)

			// Continue to next handler
			return next(c)
		}
	}
}
