package services

import (
	"context"
	"errors"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
)

// UserService implements ports.UserService interface and provides access to the user repository
type UserService struct {
	repo ports.UserRepository
}

// NewUserService creates a new user services instance
func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetByID gets a user by ID
func (us *UserService) GetByID(ctx context.Context, id int) (*entities.User, error) {
	user, err := us.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

// Register creates a new user
func (us *UserService) Register(ctx context.Context, user *entities.User) (*entities.User, error) {
	created, err := us.repo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrUserUsernameAlreadyExists) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}
	return created, nil
}
