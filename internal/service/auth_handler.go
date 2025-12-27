package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	AuthService *AuthService
}

// NewAuthHandler creates a new auth handler instance
func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

// GoogleLogin redirects the user to Google OAuth login page
func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	url := h.AuthService.GetGoogleLoginURL()
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the OAuth callback from Google
func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	// Get authorization code from query parameters
	code := c.QueryParam("code")
	if code == "" {
		return c.String(http.StatusBadRequest, "Authorization code not found")
	}

	// Exchange code for user information and tokens
	user, err := h.AuthService.HandleGoogleCallback(code)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to authenticate: "+err.Error())
	}

	// Generate JWT session token
	accessToken, err := h.AuthService.GenerateSessionToken(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to generate session token: "+err.Error())
	}

	// Set JWT as HTTP-only cookie
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   h.AuthService.JWTExpiry * 3600, // Convert hours to seconds
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	// Redirect to contacts page
	return c.Redirect(http.StatusFound, "/contacts")
}

// Logout clears the authentication cookie
func (h *AuthHandler) Logout(c echo.Context) error {
	// Clear the auth cookie
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Delete cookie
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/login")
}
