//go:build !integration

package services_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/services"
	"go-starter/internal/domain/utils"
	"testing"
)

func TestUserService_Register(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	userToCreate := &entities.User{
		Username: "example",
		Password: "secret123",
	}

	createdUser, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	tests := map[string]struct {
		input       *entities.User
		expectedErr error
	}{
		"register valid user": {
			input: &entities.User{
				Username: "success",
				Password: "secret123",
			},
			expectedErr: nil,
		},
		"register user with conflicting username": {
			input: &entities.User{
				Username: createdUser.Username,
				Password: userToCreate.Password,
			},
			expectedErr: domain.ErrUsernameAlreadyTaken,
		},
		"register user without username": {
			input: &entities.User{
				Password: userToCreate.Password,
			},
			expectedErr: domain.ErrUsernameRequired,
		},
		"register user without password": {
			input: &entities.User{
				Username: createdUser.Username,
			},
			expectedErr: domain.ErrPasswordRequired,
		},
		"register user with short password": {
			input: &entities.User{
				Username: createdUser.Username,
				Password: "short",
			},
			expectedErr: domain.ErrPasswordTooShort,
		},
		"register user with short username": {
			input: &entities.User{
				Username: "abc",
				Password: userToCreate.Password,
			},
			expectedErr: domain.ErrUsernameTooShort,
		},
		"register user with long username": {
			input: &entities.User{
				Username: "this_is_a_really_long_username",
				Password: userToCreate.Password,
			},
			expectedErr: domain.ErrUsernameTooLong,
		},
		"register user with invalid username": {
			input: &entities.User{
				Username: " invalid ",
				Password: userToCreate.Password,
			},
			expectedErr: domain.ErrUsernameInvalid,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result, err := builder.UserService.Register(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err == nil && result.Username != tt.input.Username {
				t.Errorf("expected username %s, got %s", tt.input.Username, result.Username)
			}
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	userToCreate := &entities.User{
		Username: "example",
		Password: "secret123",
	}

	createdUser, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	tests := map[string]struct {
		input       entities.UserID
		expectedErr error
	}{
		"get valid user by id": {
			input:       createdUser.ID,
			expectedErr: nil,
		},
		"get non-existing user by id": {
			input:       entities.UserID(uuid.New()),
			expectedErr: domain.ErrUserNotFound,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := builder.UserService.GetByID(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUserService_GetByID_Cache(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	userToCreate := &entities.User{
		Username: "example",
		Password: "secret123",
	}
	createdUser, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	_, err = builder.UserService.GetByID(ctx, createdUser.ID)
	if err != nil {
		t.Fatalf("error while fetching user: %v", err)
	}

	// Act & Assert
	cachedUser, err := builder.CacheService.Get(ctx, utils.GenerateCacheKey(services.UserCachePrefix, createdUser.ID))
	if err != nil {
		t.Errorf("error while getting user from cache: %v", err)
	}

	var deserializedUser entities.User
	err = utils.Deserialize(cachedUser, &deserializedUser)
	if err != nil {
		t.Errorf("error while deserializing user: %v", err)
	}

	if deserializedUser.ID != createdUser.ID {
		t.Errorf("deserialized user does not match cache")
	}
}

func TestUserService_GetByUsername(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	userToCreate := &entities.User{
		Username: "example",
		Password: "secret123",
	}

	createdUser, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	tests := map[string]struct {
		input       string
		expectedErr error
	}{
		"get valid user by username": {
			input:       createdUser.Username,
			expectedErr: nil,
		},
		"get non-existing user by username": {
			input:       "non-existing",
			expectedErr: domain.ErrUserNotFound,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result, err := builder.UserService.GetByUsername(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err == nil && result.Username != tt.input {
				t.Errorf("expected username %s, got %s", tt.input, result.Username)
			}
		})
	}
}

func TestUserService_UpdatePassword(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()
	userToCreate := &entities.User{
		Username: "example",
		Password: "secret123",
	}
	user, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	validPassword := "secret123"
	shortPassword := "short"
	notMatchingPassword := "not-matching"

	tests := map[string]struct {
		input       entities.UpdateUserParams
		expectedErr error
	}{
		"update password successfully": {
			input: entities.UpdateUserParams{
				Password:             &validPassword,
				PasswordConfirmation: &validPassword,
			},
			expectedErr: nil,
		},
		"update with short password": {
			input: entities.UpdateUserParams{
				Password:             &shortPassword,
				PasswordConfirmation: &shortPassword,
			},
			expectedErr: domain.ErrPasswordTooShort,
		},
		"update with non-matching password": {
			input: entities.UpdateUserParams{
				Password:             &validPassword,
				PasswordConfirmation: &notMatchingPassword,
			},
			expectedErr: domain.ErrPasswordsNotMatch,
		},
		"update with missing password": {
			input: entities.UpdateUserParams{
				PasswordConfirmation: &validPassword,
			},
			expectedErr: domain.ErrPasswordRequired,
		},
		"update with missing password confirmation": {
			input: entities.UpdateUserParams{
				Password: &validPassword,
			},
			expectedErr: domain.ErrPasswordConfirmationRequired,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := builder.UserService.UpdatePassword(ctx, user.ID, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
