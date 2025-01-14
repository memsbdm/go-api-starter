//go:build !integration

package services_test

import (
	"errors"
	"github.com/google/uuid"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"testing"
	"time"
)

var refreshTokenTimeToExpire = 2 * time.Hour
var tokenTimeToExpire = 20 * time.Minute

func TestTokenService_ValidateToken(t *testing.T) {
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
			advance:     tokenTimeToExpire,
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
			_, err = builder.TokenService.GetTokenPayload(accessToken)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestTokenService_ValidateRefreshToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		advance     time.Duration
		expectedErr error
	}{
		{
			name:        "Valid refresh token",
			advance:     0,
			expectedErr: nil,
		},
		{
			name:        "Expired refresh token",
			advance:     refreshTokenTimeToExpire,
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

			refreshToken, err := builder.TokenService.GenerateRefreshToken(user)
			if err != nil {
				t.Fatalf("failed to login: %v", err)
			}

			builder.TimeGenerator.Advance(tt.advance)

			// Act & Assert
			_, err = builder.TokenService.ValidateRefreshToken(refreshToken)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
