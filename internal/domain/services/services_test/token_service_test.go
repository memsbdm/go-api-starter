//go:build !integration

package services_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"testing"
	"time"
)

var refreshTokenExpirationDuration = 2 * time.Hour
var accessTokenExpirationDuration = 20 * time.Minute

func TestTokenService_ValidateAndParseAccessToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		advance     time.Duration
		expectedErr error
	}{
		{
			name:        "Valid token",
			advance:     0,
			expectedErr: nil,
		},
		{
			name:        "Expired token",
			advance:     accessTokenExpirationDuration,
			expectedErr: domain.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			timeGenerator := timegen.NewFakeTimeGenerator(time.Now())
			builder := NewTestBuilder().WithTimeGenerator(timeGenerator).Build()

			user := &entities.User{
				ID: entities.UserID(uuid.New()),
			}

			accessToken, err := builder.TokenService.GenerateAccessToken(user)
			if err != nil {
				t.Fatalf("failed to generated access token: %v", err)
			}

			builder.TimeGenerator.Advance(tt.advance)

			// Act & Assert
			claims, err := builder.TokenService.ValidateAndParseAccessToken(accessToken)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil && claims.Subject != user.ID {
				t.Errorf("expected user id %s, got %s", user.ID, claims.Subject)
			}
		})
	}
}

func TestTokenService_ValidateAndParseRefreshToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		advance     time.Duration
		expectedErr error
	}{
		{
			name:        "Validate and parse valid refresh token",
			advance:     0,
			expectedErr: nil,
		},
		{
			name:        "Validate and parse expired refresh token",
			advance:     refreshTokenExpirationDuration,
			expectedErr: domain.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctx := context.Background()
			timeGenerator := timegen.NewFakeTimeGenerator(time.Now())
			builder := NewTestBuilder().WithTimeGenerator(timeGenerator).Build()
			user := &entities.User{
				ID: entities.UserID(uuid.New()),
			}
			refreshToken, err := builder.TokenService.GenerateRefreshToken(ctx, user.ID)
			if err != nil {
				t.Fatalf("failed to generated refresh token: %v", err)
			}
			builder.TimeGenerator.Advance(tt.advance)

			// Act & Assert
			claims, err := builder.TokenService.ValidateAndParseRefreshToken(ctx, refreshToken)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err == nil && claims.Subject != user.ID {
				t.Errorf("expected user id %s, got %s", user.ID, claims.Subject)
			}
		})
	}
}

func TestTokenService_RevokeRefreshToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		advance     time.Duration
		expectedErr error
	}{
		{
			name:        "Revoke valid refresh token",
			advance:     0,
			expectedErr: nil,
		},
		{
			name:        "Revoke expired refresh token",
			advance:     accessTokenExpirationDuration,
			expectedErr: domain.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctx := context.Background()
			timeGenerator := timegen.NewFakeTimeGenerator(time.Now())
			builder := NewTestBuilder().WithTimeGenerator(timeGenerator).Build()
			user := &entities.User{
				ID: entities.UserID(uuid.New()),
			}

			refreshToken, err := builder.TokenService.GenerateRefreshToken(ctx, user.ID)
			if err != nil {
				t.Fatalf("failed to generated refresh token: %v", err)
			}
			builder.TimeGenerator.Advance(tt.advance)

			err = builder.TokenService.RevokeRefreshToken(ctx, refreshToken)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
