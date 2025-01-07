//go:build !integration

package services_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go-starter/internal/adapters/storage/postgres/repositories/mocks"
	"go-starter/internal/adapters/storage/redis"
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
		cacheRepo := redis.NewMock()
		cacheService := services.NewCacheService(cacheRepo)
		userRepo := mocks.MockUserRepository()
		userService = services.NewUserService(userRepo, cacheService)
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
		expectedErr error
		setupMock   func() entities.UserID
	}{
		{
			name:        "get user by id successfully",
			expectedErr: nil,
			setupMock: func() entities.UserID {
				user, _ := userService.Register(context.Background(), &entities.User{
					Username: "example",
				})
				return user.ID
			},
		},
		{
			name:        "get non-existing user by id",
			expectedErr: domain.ErrUserNotFound,
			setupMock: func() entities.UserID {
				return entities.UserID(uuid.New())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputID := tt.setupMock()
			_, err := userService.GetByID(context.Background(), inputID)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
