package ports

import (
	"context"
	"go-starter/internal/domain/entities"
	"io"
)

// UserService is an interface for interacting with user-related business logic.
type UserService interface {
	// GetByID retrieves a user by their unique identifier.
	// Returns the user entity if found or an error if not found or any other issue occurs.
	GetByID(ctx context.Context, id entities.UserID) (*entities.User, error)

	// GetByUsername retrieves a user by their username.
	// Returns the user entity if found or an error if not found or any other issue occurs.
	GetByUsername(ctx context.Context, username string) (*entities.User, error)

	// GetIDByVerifiedEmail retrieves a user ID by their verified email.
	// Returns the user ID if found or an error if not found or any other issue occurs.
	GetIDByVerifiedEmail(ctx context.Context, email string) (entities.UserID, error)

	// Register creates a new user account in the system.
	// Returns the created user or an error if the registration fails (e.g., due to validation issues).
	Register(ctx context.Context, user *entities.User) (*entities.User, error)

	// UpdatePassword updates a user password.
	// Returns an error if the update fails (e.g., due to validation issues).
	UpdatePassword(ctx context.Context, userID entities.UserID, params entities.UpdateUserParams) error

	// VerifyEmail verifies a user email.
	// Returns an error if the verification fails.
	VerifyEmail(ctx context.Context, token string) error

	// ResendEmailVerification resends a user email verification email.
	// Returns an error if the resend fails.
	ResendEmailVerification(ctx context.Context, userID entities.UserID) error

	// UpdateAvatar updates a user avatar.
	// Returns an error if the update fails.
	UpdateAvatar(ctx context.Context, userID entities.UserID, filename string, file io.Reader) (string, error)

	// DeleteAvatar deletes a user avatar.
	// Returns an error if the deletion fails.
	DeleteAvatar(ctx context.Context, userID entities.UserID) error
}

// UserRepository is an interface for interacting with user-related data.
type UserRepository interface {
	// GetByID selects a user by their unique identifier from the database.
	// Returns the user entity if found or an error if not found or any other issue occurs.
	GetByID(ctx context.Context, id entities.UserID) (*entities.User, error)

	// GetByUsername selects a user by their username from the database.
	// Returns the user entity if found or an error if not found or any other issue occurs.
	GetByUsername(ctx context.Context, username string) (*entities.User, error)

	// GetIDByVerifiedEmail returns the user ID for a verified email.
	// Returns an error if the user is not found or any other issue occurs.
	GetIDByVerifiedEmail(ctx context.Context, email string) (entities.UserID, error)

	// CheckEmailAvailability checks if an email is available for registration.
	// Returns an error if the email is already taken.
	CheckEmailAvailability(ctx context.Context, email string) error

	// Create inserts a new user into the database.
	// Returns the created user or an error if the operation fails (e.g., due to a database constraint violation).
	Create(ctx context.Context, user *entities.User) (*entities.User, error)

	// UpdatePassword updates a user password.
	// Returns an error if the update fails (e.g., due to validation issues).
	UpdatePassword(ctx context.Context, userID entities.UserID, newPassword string) error

	// VerifyEmail updates the email verification status of a user.
	// Returns the updated user or an error if the verification fails.
	VerifyEmail(ctx context.Context, userID entities.UserID) (*entities.User, error)

	// UpdateAvatar updates a user avatar.
	// Returns an error if the update fails.
	UpdateAvatar(ctx context.Context, userID entities.UserID, avatarURL string) error

	// DeleteAvatar deletes a user avatar.
	// Returns an error if the deletion fails.
	DeleteAvatar(ctx context.Context, userID entities.UserID) error
}
