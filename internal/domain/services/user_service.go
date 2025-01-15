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

// UserService implements ports.UserService interface and provides access to the user repository.
type UserService struct {
	repo  ports.UserRepository
	cache ports.CacheService
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo ports.UserRepository, cache ports.CacheService) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

// UserCachePrefix is the prefix for caching users.
const UserCachePrefix = "user"

// GetByID retrieves a user by their unique identifier.
// Returns the user entity if found or an error if not found or any other issue occurs.
func (us *UserService) GetByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	var user *entities.User
	cacheKey := utils.GenerateCacheKey(UserCachePrefix, id)
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

// GetByUsername retrieves a user by their username.
// Returns the user entity if found or an error if not found or any other issue occurs.
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

// Register creates a new user account in the system.
// Returns the created user or an error if the registration fails (e.g., due to validation issues).
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

// UpdatePassword updates a user password.
// Returns an error if the update fails (e.g., due to validation issues).
func (us *UserService) UpdatePassword(ctx context.Context, userID entities.UserID, params entities.UpdateUserParams) error {
	if params.Password == nil {
		return domain.ErrPasswordRequired
	}
	if params.PasswordConfirmation == nil {
		return domain.ErrPasswordConfirmationRequired
	}
	if *params.Password != *params.PasswordConfirmation {
		return domain.ErrPasswordsNotMatch
	}
	if len(*params.Password) < domain.PasswordMinLength {
		return domain.ErrPasswordTooShort
	}
	hashedPassword, err := utils.HashPassword(*params.Password)
	if err != nil {
		return domain.ErrInternal
	}
	err = us.repo.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}
