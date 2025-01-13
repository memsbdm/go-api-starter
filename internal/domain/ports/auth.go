package ports

import (
	"context"
	"go-starter/internal/domain/entities"
)

// AuthService is an interface for interacting with authentication
type AuthService interface {
	// Login authenticates a user.
	// Returns an access token, a refresh token or an error if login fails.
	Login(ctx context.Context, username, password string) (string, string, error)
	// RefreshToken generates a new access token and a new refresh token. It returns an error if the token is invalid or expired
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	// Register registers a new user
	Register(ctx context.Context, user *entities.User) (*entities.User, error)
}
