package ports

import (
	"context"
	"go-starter/internal/domain/entities"
)

// AuthService is an interface for interacting with authentication
type AuthService interface {
	// Login authenticates a user.
	// Returns an access token, a refresh token or an error if login fails
	Login(ctx context.Context, username, password string) (string, string, error)
}

// TokenService is an interface for interacting with token-related business logic
type TokenService interface {
	// CreateToken creates a new token for a given user
	CreateToken(user *entities.User, isRefreshToken bool) (string, error)
	// ValidateToken validates an auth token and returns associated token payload
	ValidateToken(tokenStr string) (*entities.TokenPayload, error)
	// CreateRefreshToken creates a new refresh token for a given user
	CreateRefreshToken(user *entities.User) (string, error)
	// ValidateRefreshToken validates a refresh token and returns associated token payload
	ValidateRefreshToken(tokenStr string) error
}
