package ports

import (
	"context"
	"go-starter/internal/domain/entities"
)

// UserService is an interface for interacting with user-related business logic.
type UserService interface {
	// GetByID retrieves a user by their unique identifier.
	// Returns the user entity if found or an error if not found or any other issue occurs.
	GetByID(ctx context.Context, id entities.UserID) (*entities.User, error)

	// GetByUsername retrieves a user by their username.
	// Returns the user entity if found or an error if not found or any other issue occurs.
	GetByUsername(ctx context.Context, username string) (*entities.User, error)

	// Register creates a new user account in the system.
	// Returns the created user or an error if the registration fails (e.g., due to validation issues).
	Register(ctx context.Context, user *entities.User) (*entities.User, error)
}

// UserRepository is an interface for interacting with user-related data.
type UserRepository interface {
	// GetByID selects a user by their unique identifier from the database.
	// Returns the user entity if found or an error if not found or any other issue occurs.
	GetByID(ctx context.Context, id entities.UserID) (*entities.User, error)

	// GetByUsername selects a user by their username from the database.
	// Returns the user entity if found or an error if not found or any other issue occurs.
	GetByUsername(ctx context.Context, username string) (*entities.User, error)

	// Create inserts a new user into the database.
	// Returns the created user or an error if the operation fails (e.g., due to a database constraint violation).
	Create(ctx context.Context, user *entities.User) (*entities.User, error)
}
