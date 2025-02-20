package mocks

import (
	"context"
	"github.com/google/uuid"
	"go-starter/internal/adapters/storage/postgres/repositories"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"sync"
)

type db struct {
	data map[entities.UserID]*entities.User
	mu   sync.RWMutex
}

// UserRepository implements the ports.UserRepository interface and provides access to the database.
type UserRepository struct {
	db db
}

// NewUserRepositoryMock creates and returns a new mock instance of a user repository.
func NewUserRepositoryMock() *UserRepository {
	return &UserRepository{
		db: db{
			data: map[entities.UserID]*entities.User{},
			mu:   sync.RWMutex{},
		},
	}
}

// GetByID selects a user by their unique identifier from the database.
// Returns the user entity if found or an error if not found or any other issue occurs.
func (ur *UserRepository) GetByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, repositories.QueryTimeoutDuration)
	defer cancel()

	ur.db.mu.Lock()
	defer ur.db.mu.Unlock()

	user, ok := ur.db.data[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

// GetByUsername selects a user by their username from the database.
// Returns the user entity if found or an error if not found or any other issue occurs.
func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, repositories.QueryTimeoutDuration)
	defer cancel()
	ur.db.mu.Lock()
	defer ur.db.mu.Unlock()
	for _, v := range ur.db.data {
		if v.Username == username {
			return v, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

// Create inserts a new user into the database.
// Returns the created user or an error if the operation fails (e.g., due to a database constraint violation).
func (ur *UserRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, repositories.QueryTimeoutDuration)
	defer cancel()

	ur.db.mu.Lock()
	defer ur.db.mu.Unlock()

	for _, v := range ur.db.data {
		if v.Username == user.Username {
			return nil, domain.ErrUsernameAlreadyTaken
		}
	}

	id := uuid.New()
	newUser := &entities.User{
		ID:       entities.UserID(id),
		Username: user.Username,
		Password: user.Password,
	}
	ur.db.data[newUser.ID] = newUser

	return newUser, nil
}

// UpdatePassword updates a user password.
// Returns an error if the update fails (e.g., due to validation issues).
func (ur *UserRepository) UpdatePassword(ctx context.Context, userID entities.UserID, newPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, repositories.QueryTimeoutDuration)
	defer cancel()

	ur.db.mu.Lock()
	defer ur.db.mu.Unlock()

	ur.db.data[userID].Username = newPassword
	return nil
}
