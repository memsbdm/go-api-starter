//go:build !integration

package services_test

import (
	"context"
	"errors"
	"go-starter/internal/adapters/mocks"
	"go-starter/internal/domain"
	"testing"
	"time"
)

func TestAuthService_Login(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	userToCreate := newValidUserToCreate()

	_, err := builder.AuthService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	type loginRequest struct {
		username string
		password string
	}

	tests := map[string]struct {
		input       *loginRequest
		expectedErr error
	}{
		"login success": {
			input: &loginRequest{
				username: userToCreate.Username,
				password: userToCreate.Password,
			},
			expectedErr: nil,
		},
		"login with bad credentials": {
			input: &loginRequest{
				username: "not-existing",
				password: "not-existing",
			},
			expectedErr: domain.ErrInvalidCredentials,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, _, err := builder.AuthService.Login(ctx, tt.input.username, tt.input.password)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestAuthService_Refresh(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		advance     time.Duration
		modifyToken func() string
		expectedErr error
	}{
		"refresh with a valid refresh token": {
			advance:     0,
			modifyToken: nil,
			expectedErr: nil,
		},
		"refresh with an expired refresh token": {
			advance:     refreshTokenExpirationDuration,
			modifyToken: nil,
			expectedErr: domain.ErrInvalidToken,
		},
		"refresh without a refresh token": {
			advance:     0,
			modifyToken: func() string { return "" },
			expectedErr: domain.ErrRefreshTokenRequired,
		},
		"refresh with an invalid refresh token": {
			advance:     0,
			modifyToken: func() string { return "invalid" },
			expectedErr: domain.ErrInvalidToken,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctx := context.Background()
			timeGenerator := mocks.NewTimeGeneratorMock(time.Now())
			builder := NewTestBuilder().WithTimeGenerator(timeGenerator).Build()

			userToCreate := newValidUserToCreate()

			_, err := builder.AuthService.Register(ctx, userToCreate)
			if err != nil {
				t.Fatalf("failed to register user: %v", err)
			}

			_, authTokens, err := builder.AuthService.Login(ctx, userToCreate.Username, userToCreate.Password)
			if err != nil {
				t.Fatalf("failed to login: %v", err)
			}

			advanceTime(t, builder.TimeGenerator, tt.advance)

			refreshToken := authTokens.RefreshToken
			if tt.modifyToken != nil {
				refreshToken = tt.modifyToken()
			}

			// Act & Assert
			_, err = builder.AuthService.Refresh(ctx, refreshToken)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestAuthService_Register_SendsEmail(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	// Act & Assert
	userToCreate := newValidUserToCreate()
	_, err := builder.AuthService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	if v, ok := builder.MailerAdapter.(interface{ SentEmailsCount() int }); ok {
		if v.SentEmailsCount() != 1 {
			t.Errorf("expected 1 email to be sent, got %d", v.SentEmailsCount())
		}
	}
}
