package ports

import (
	"context"
	"go-starter/internal/domain/entities"
)

// AuthService is an interface for interacting with authentication operations.
type AuthService interface {
	// Login logs in a user in the system.
	// Returns the user entity and the access token and an error if the login fails
	// (e.g., due to invalid credentials or validation issues).
	Login(ctx context.Context, username, password string) (*entities.User, string, error)

	// Register registers a new user in the system.
	// Returns the created user entity and an error if the registration fails
	// (e.g., due to username already existing or validation issues).
	Register(ctx context.Context, user *entities.User) (*entities.User, error)

	// Logout logs out a user from the system.
	// Returns an error if the logout fails.
	Logout(ctx context.Context, accessToken string) error
}
