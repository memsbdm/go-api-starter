package services

import (
	"context"
	"errors"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/utils"
	"time"
)

// UserService implements ports.UserService interface and provides access to the user repository
type UserService struct {
	repo  ports.UserRepository
	cache ports.CacheService
}

// NewUserService creates a new user services instance
func NewUserService(repo ports.UserRepository, cache ports.CacheService) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

// GetByID gets a user by ID
func (us *UserService) GetByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	var user *entities.User
	cacheKey := utils.GenerateCacheKey("user", id)
	cachedUser, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err := utils.Deserialize(cachedUser, &user)
		if err != nil {
			return nil, domain.ErrInternal
		}
		return user, nil
	}

	user, err = us.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	userSerialized, err := utils.Serialize(user)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, time.Hour)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return user, nil
}

// GetByUsername gets a user by username
func (us *UserService) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	user, err := us.repo.GetByUsername(ctx, username)
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
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, domain.ErrInternal
	}

	userToCreate := &entities.User{
		Username: user.Username,
		Password: hashedPassword,
	}

	created, err := us.repo.Create(ctx, userToCreate)
	if err != nil {
		if errors.Is(err, domain.ErrUserUsernameAlreadyExists) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return created, nil
}
