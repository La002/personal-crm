package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/La002/personal-crm/pkg/entity"
	jwtutil "github.com/La002/personal-crm/pkg/jwt"
	"github.com/La002/personal-crm/pkg/repository"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleUserInfo represents the user information returned by Google OAuth
type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// AuthService handles authentication operations
type AuthService struct {
	UserRepo     repository.UserDao
	OAuth2Config *oauth2.Config
	JWTSecret    string
	JWTExpiry    int
}

// NewAuthService creates a new auth service instance
func NewAuthService(
	userRepo repository.UserDao,
	clientID string,
	clientSecret string,
	redirectURL string,
	jwtSecret string,
	jwtExpiry int,
) *AuthService {
	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/calendar.events",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthService{
		UserRepo:     userRepo,
		OAuth2Config: oauth2Config,
		JWTSecret:    jwtSecret,
		JWTExpiry:    jwtExpiry,
	}
}

// GetGoogleLoginURL generates the Google OAuth login URL
func (s *AuthService) GetGoogleLoginURL() string {
	return s.OAuth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

// HandleGoogleCallback exchanges the authorization code for tokens and creates/updates user
func (s *AuthService) HandleGoogleCallback(code string) (*entity.User, error) {
	// Exchange code for token
	ctx := context.Background()
	token, err := s.OAuth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from Google
	client := s.OAuth2Config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var googleUser GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Check if user already exists by GoogleID
	user, err := s.UserRepo.GetUserByGoogleID(googleUser.ID)
	if err != nil {
		// User doesn't exist by GoogleID, check by email
		user, emailErr := s.UserRepo.GetUserByEmail(googleUser.Email)
		if emailErr != nil {
			// User doesn't exist at all, create new user
			user = entity.User{
				GoogleId:     googleUser.ID,
				Email:        googleUser.Email,
				Picture:      googleUser.Picture,
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				TokenExpiry:  token.Expiry,
			}

			if err := s.UserRepo.CreateUser(&user); err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
			return &user, nil
		}
		// User exists by email, update GoogleID and tokens
		user.GoogleId = googleUser.ID
		user.AccessToken = token.AccessToken
		user.RefreshToken = token.RefreshToken
		user.TokenExpiry = token.Expiry
		user.Picture = googleUser.Picture

		if err := s.UserRepo.UpdateUser(&user); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		// User exists by GoogleID, update tokens
		user.AccessToken = token.AccessToken
		user.RefreshToken = token.RefreshToken
		user.TokenExpiry = token.Expiry
		user.Picture = googleUser.Picture

		if err := s.UserRepo.UpdateUser(&user); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	return &user, nil
}

// GenerateSessionToken creates a JWT token for the user
func (s *AuthService) GenerateSessionToken(user *entity.User) (string, error) {
	return jwtutil.GenerateToken(user.ID, user.Email, s.JWTSecret, s.JWTExpiry)
}
