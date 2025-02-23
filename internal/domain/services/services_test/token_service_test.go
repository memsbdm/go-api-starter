//go:build !integration

package services_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go-starter/internal/adapters/mocks"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"testing"
	"time"
)

var refreshTokenExpirationDuration = 2 * time.Hour
var accessTokenExpirationDuration = 20 * time.Minute

func TestTokenService_ValidateAndParse(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input       entities.TokenType
		advance     time.Duration
		expectedErr error
	}{
		"validate and parse valid access token": {
			input:       entities.AccessToken,
			advance:     0,
			expectedErr: nil,
		},
		"validate and parse expired access token": {
			input:       entities.AccessToken,
			advance:     accessTokenExpirationDuration,
			expectedErr: domain.ErrInvalidToken,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			timeGenerator := mocks.NewTimeGeneratorMock(time.Now())
			builder := NewTestBuilder().WithTimeGenerator(timeGenerator).Build()

			user := &entities.User{
				ID: entities.UserID(uuid.New()),
			}

			token, err := builder.TokenService.Generate(tt.input, user)
			if err != nil {
				t.Fatalf("failed to generated access token: %v", err)
			}

			advanceTime(t, builder.TimeGenerator, tt.advance)

			// Act & Assert
			claims, err := builder.TokenService.ValidateAndParse(tt.input, token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil && claims.Subject != user.ID {
				t.Errorf("expected user id %s, got %s", user.ID, claims.Subject)
			}
		})
	}
}

func TestTokenService_ValidateAndParseWithCache(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input       entities.TokenType
		advance     time.Duration
		expectedErr error
	}{
		"validate and parse valid refresh token": {
			input:       entities.RefreshToken,
			advance:     0,
			expectedErr: nil,
		},
		"validate and parse expired refresh token": {
			input:       entities.RefreshToken,
			advance:     refreshTokenExpirationDuration,
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
			user := &entities.User{
				ID: entities.UserID(uuid.New()),
			}
			token, err := builder.TokenService.GenerateTokenWithCache(ctx, tt.input, user)
			if err != nil {
				t.Fatalf("failed to generated refresh token: %v", err)
			}
			advanceTime(t, builder.TimeGenerator, tt.advance)

			// Act & Assert
			claims, err := builder.TokenService.ValidateAndParseWithCache(ctx, tt.input, token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err == nil && claims.Subject != user.ID {
				t.Errorf("expected user id %s, got %s", user.ID, claims.Subject)
			}
		})
	}
}

func TestTokenService_RevokeTokenFromCache(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input       entities.TokenType
		advance     time.Duration
		expectedErr error
	}{
		"revoke valid refresh token": {
			input:       entities.RefreshToken,
			advance:     0,
			expectedErr: nil,
		},
		"revoke expired refresh token": {
			input:       entities.RefreshToken,
			advance:     refreshTokenExpirationDuration,
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
			user := &entities.User{
				ID: entities.UserID(uuid.New()),
			}

			token, err := builder.TokenService.GenerateTokenWithCache(ctx, tt.input, user)
			if err != nil {
				t.Fatalf("failed to generated refresh token: %v", err)
			}
			advanceTime(t, builder.TimeGenerator, tt.advance)

			err = builder.TokenService.RevokeTokenFromCache(ctx, tt.input, token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
