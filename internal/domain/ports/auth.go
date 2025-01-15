package ports

import (
	"context"
	"go-starter/internal/domain/entities"
)

// AuthService is an interface for interacting with authentication
type AuthService interface {
	// Login authenticates a user.
	// Returns an access token, a refresh token or an error if login fails.
	Login(ctx context.Context, username, password string) (accessToken string, refreshToken string, err error)
	// Register registers a new user
	Register(ctx context.Context, user *entities.User) (*entities.User, error)
	Refresh(ctx context.Context, previousRefreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
}
