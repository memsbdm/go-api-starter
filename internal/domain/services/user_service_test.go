package services_test

import (
	"context"
	"errors"
	"go-starter/internal/adapters/storage/postgres/repositories/mocks"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/services"
	"sync"
	"testing"
)

var (
	once        sync.Once
	userService *services.UserService
)

func setup() {
	once.Do(func() {
		userRepo := mocks.MockUserRepository()
		userService = services.NewUserService(userRepo)
	})
}

func TestUserService_Register(t *testing.T) {
	t.Parallel()
	setup()

	tests := []struct {
		name        string
		input       *entities.User
		expectedErr error
		setupMock   func()
	}{
		{
			name: "create user successfully",
			input: &entities.User{
				Username: "success",
			},
			expectedErr: nil,
			setupMock:   func() {},
		},
		{
			name: "create user with conflicting username",
			input: &entities.User{
				Username: "conflict",
			},
			expectedErr: domain.ErrUserUsernameAlreadyExists,
			setupMock: func() {
				_, _ = userService.Register(context.Background(), &entities.User{
					Username: "conflict",
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			_, err := userService.Register(context.Background(), tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	t.Parallel()
	setup()

	tests := []struct {
		name        string
		input       int
		expectedErr error
		setupMock   func()
	}{
		{
			name:        "get user by id successfully",
			input:       0,
			expectedErr: nil,
			setupMock: func() {
				_, _ = userService.Register(context.Background(), &entities.User{
					Username: "not empty",
				})
			},
		},
		{
			name:        "get non-existing user by id",
			input:       -1,
			expectedErr: domain.ErrUserNotFound,
			setupMock:   func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			_, err := userService.GetByID(context.Background(), tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
