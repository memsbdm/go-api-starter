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

	// SendPasswordResetEmail sends a password reset email to the user.
	// Returns an error if the email fails to send.
	SendPasswordResetEmail(ctx context.Context, email string) error

	// VerifyPasswordResetToken verifies a password reset token.
	// Returns an error if the token is invalid.
	VerifyPasswordResetToken(ctx context.Context, token string) error
}
