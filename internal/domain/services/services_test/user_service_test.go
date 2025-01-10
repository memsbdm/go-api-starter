// //go:build !integration

package services_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"testing"
)

func TestUserService_Register(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()

	userToCreate := &entities.User{
		Username: "register_init",
		Password: "secret123",
	}

	createdUser, err := userService.Register(ctx, userToCreate)
	if err != nil {
		t.Errorf("error while registering user: %v", err)
	}

	tests := []struct {
		name        string
		input       *entities.User
		expectedErr error
	}{
		{
			name: "create user successfully",
			input: &entities.User{
				Username: "register",
				Password: "secret123",
			},
			expectedErr: nil,
		},
		{
			name: "create user with conflicting username",
			input: &entities.User{
				Username: createdUser.Username,
				Password: "secret123",
			},
			expectedErr: domain.ErrUserUsernameAlreadyExists,
		},
	}

	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := userService.Register(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err == nil && result.Username != tt.input.Username {
				t.Fatalf("expected username %s, got %s", tt.input.Username, result.Username)
			}
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()

	userToCreate := &entities.User{
		Username: "get_by_id_init",
		Password: "secret123",
	}

	createdUser, err := userService.Register(ctx, userToCreate)
	if err != nil {
		t.Errorf("error while registering user: %v", err)
	}

	tests := []struct {
		name        string
		input       entities.UserID
		expectedErr error
	}{
		{
			name:        "get user by id successfully",
			input:       createdUser.ID,
			expectedErr: nil,
		},
		{
			name:        "get non-existing user by id",
			input:       entities.UserID(uuid.New()),
			expectedErr: domain.ErrUserNotFound,
		},
	}

	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := userService.GetByID(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUserService_GetByUsername(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	userToCreate := &entities.User{
		Username: "get_by_username_init",
		Password: "secret123",
	}

	createdUser, err := userService.Register(ctx, userToCreate)
	if err != nil {
		t.Errorf("error while registering user: %v", err)
	}

	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "get user by username successfully",
			input:       createdUser.Username,
			expectedErr: nil,
		},
		{
			name:        "get non-existing user by username",
			input:       "non-existing",
			expectedErr: domain.ErrUserNotFound,
		},
	}

	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := userService.GetByUsername(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err == nil && result.Username != tt.input {
				t.Fatalf("expected username %s, got %s", tt.input, result.Username)
			}
		})
	}
}
