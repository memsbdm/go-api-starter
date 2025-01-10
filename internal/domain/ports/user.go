package ports

import (
	"context"
	"go-starter/internal/domain/entities"
)

// UserService is an interface for interacting with user-related business logic
type UserService interface {
	// GetByID returns a user by id
	GetByID(ctx context.Context, id entities.UserID) (*entities.User, error)
	// GetByUsername returns a user by id
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	// Register registers a new user
	Register(ctx context.Context, user *entities.User) (*entities.User, error)
}

// UserRepository is an interface for interacting with user-related data
type UserRepository interface {
	// GetByID selects a user by id
	GetByID(ctx context.Context, id entities.UserID) (*entities.User, error)
	// GetByUsername selects a user by id
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	// Create inserts a new user into the database
	Create(ctx context.Context, user *entities.User) (*entities.User, error)
}
