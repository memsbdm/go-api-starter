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
	mu   sync.Mutex
}

// UserRepository implements ports.UserRepository interface and provides access to the database
type UserRepository struct {
	db db
}

// MockUserRepository creates a new mock for a user repository instance
func MockUserRepository() *UserRepository {
	return &UserRepository{
		db: db{
			data: map[entities.UserID]*entities.User{},
			mu:   sync.Mutex{},
		},
	}
}

// GetByID gets a user by ID from the database
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

// Create creates a new user in the database
func (ur *UserRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, repositories.QueryTimeoutDuration)
	defer cancel()

	ur.db.mu.Lock()
	defer ur.db.mu.Unlock()

	for _, v := range ur.db.data {
		if v.Username == user.Username {
			return nil, domain.ErrUserUsernameAlreadyExists
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

// GetByUsername gets a user by username from the database
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