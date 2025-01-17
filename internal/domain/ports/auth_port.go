package ports

import (
	"context"
	"go-starter/internal/domain/entities"
)

// AuthService is an interface for interacting with authentication operations.
type AuthService interface {
	// Login authenticates a user.
	// Returns auth tokens upon successful authentication,
	// or an error if the login fails (e.g., due to incorrect credentials).
	Login(ctx context.Context, username, password string) (*entities.User, *entities.AuthTokens, error)

	// Register registers a new user in the system.
	// Returns the created user entity and an error if the registration fails
	// (e.g., due to username already existing or validation issues).
	Register(ctx context.Context, user *entities.User) (*entities.User, error)

	// Refresh generates new access and refresh tokens using the previous refresh token.
	// Returns the new auth tokens, and an error if the refresh fails
	// (e.g., if the previous refresh token is invalid or expired).
	Refresh(ctx context.Context, previousRefreshToken string) (*entities.AuthTokens, error)

	// Logout invalidates the specified refresh token, effectively logging the user out.
	// Returns an error if the logout operation fails (e.g., if the refresh token is not found).
	Logout(ctx context.Context, refreshToken string) error
}
