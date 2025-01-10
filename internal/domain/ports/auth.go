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
	// ValidateToken validates an auth token (JWT, session ID, etc.) and returns associated token payload
	ValidateToken(tokenStr string) (*entities.TokenPayload, error)
}
