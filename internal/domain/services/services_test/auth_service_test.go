package services_test

import (
	"context"
	"errors"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"testing"
)

func TestAuthService_Login(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	userToCreate := &entities.User{
		Username: "login_init",
		Password: "secret123",
	}

	_, err := userService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	type loginRequest struct {
		username string
		password string
	}

	tests := []struct {
		name        string
		input       *loginRequest
		expectedErr error
	}{
		{
			name: "Login Success",
			input: &loginRequest{
				username: userToCreate.Username,
				password: userToCreate.Password,
			},
			expectedErr: nil,
		},
		{
			name: "Login Error",
			input: &loginRequest{
				username: "not-existing",
				password: "not-existing",
			},
			expectedErr: domain.ErrInvalidCredentials,
		},
	}

	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, _, err := authService.Login(ctx, tt.input.username, tt.input.password)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
