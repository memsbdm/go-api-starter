package mocks

import (
	"context"
	"go-starter/internal/adapters/storage/postgres/repositories"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"sync"
)

type db struct {
	data map[int]*entities.User
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
			data: map[int]*entities.User{},
			mu:   sync.Mutex{},
		},
	}
}

// GetByID gets a user by ID from the database
func (ur *UserRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
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

	newUser := &entities.User{
		ID:       len(ur.db.data),
		Username: user.Username,
	}
	ur.db.data[newUser.ID] = newUser

	return newUser, nil
}
