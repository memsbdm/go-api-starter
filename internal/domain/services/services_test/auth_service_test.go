//go:build !integration

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
	builder := NewTestBuilder().Build()

	userToCreate := &entities.User{
		Username: "example",
		Password: "secret123",
	}

	_, err := builder.UserService.Register(ctx, userToCreate)
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
			_, _, err := builder.AuthService.Login(ctx, tt.input.username, tt.input.password)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()
	userToCreate := &entities.User{
		Username: "example",
		Password: "secret123",
	}

	_, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	_, refreshToken, err := builder.AuthService.Login(ctx, userToCreate.Username, userToCreate.Password)
	if err != nil {
		t.Fatalf("failed to login: %v", err)
	}

	tests := []struct {
		name        string
		token       string
		expectedErr error
	}{
		{
			name:        "Refresh Token Success",
			token:       refreshToken,
			expectedErr: nil,
		},
		{
			name:        "Refresh Token Error",
			token:       refreshToken[:len(refreshToken)-1] + "1",
			expectedErr: domain.ErrInvalidToken,
		},
	}

	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, _, err := builder.AuthService.RefreshToken(ctx, tt.token)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
