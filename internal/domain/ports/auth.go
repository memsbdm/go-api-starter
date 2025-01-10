package ports

import (
	"context"
	"go-starter/internal/domain/entities"
)

// AuthService is an interface for interacting with authentication
type AuthService interface {
	// Login authenticates a user.
	// Returns a token string pointer (could be JWT, session ID, etc.) or an error if login fails
	Login(ctx context.Context, username, password string) (*string, error)
}

// TokenService is an interface for interacting with token-related business logic
type TokenService interface {
	// CreateToken creates a new token for a given user
	CreateToken(user *entities.User) (string, error)
	// ValidateToken validates an auth token (JWT, session ID, etc.) and returns associated token payload
	ValidateToken(tokenStr string) (*entities.TokenPayload, error)
}
