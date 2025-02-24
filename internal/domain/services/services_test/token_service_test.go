//go:build !integration

package services_test

import (
	"context"
	"errors"
	"fmt"
	"go-starter/internal/adapters/mocks"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/utils"
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	refreshTokenExpirationDuration           = 2 * time.Hour
	accessTokenExpirationDuration            = 20 * time.Minute
	emailVerificationTokenExpirationDuration = 24 * time.Hour
)

func TestTokenService_ValidateJWT(t *testing.T) {
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

			token, err := builder.TokenService.GenerateJWT(tt.input, user)
			if err != nil {
				t.Fatalf("failed to generated access token: %v", err)
			}

			advanceTime(t, builder.TimeGenerator, tt.advance)

			// Act & Assert
			claims, err := builder.TokenService.ValidateJWT(tt.input, token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil && claims.Subject != user.ID {
				t.Errorf("expected user id %s, got %s", user.ID, claims.Subject)
			}
		})
	}
}

func TestTokenService_VerifyCachedJWT(t *testing.T) {
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
			token, err := builder.TokenService.CreateAndCacheJWT(ctx, tt.input, user)
			if err != nil {
				t.Fatalf("failed to generated refresh token: %v", err)
			}
			advanceTime(t, builder.TimeGenerator, tt.advance)

			// Act & Assert
			claims, err := builder.TokenService.VerifyCachedJWT(ctx, tt.input, token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err == nil && claims.Subject != user.ID {
				t.Errorf("expected user id %s, got %s", user.ID, claims.Subject)
			}
		})
	}
}

func TestTokenService_RevokeJWT(t *testing.T) {
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

			// Act & Assert
			token, err := builder.TokenService.CreateAndCacheJWT(ctx, tt.input, user)
			if err != nil {
				t.Fatalf("failed to generated refresh token: %v", err)
			}
			advanceTime(t, builder.TimeGenerator, tt.advance)

			err = builder.TokenService.RevokeJWT(ctx, tt.input, token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			_, err = builder.TokenService.VerifyCachedJWT(ctx, tt.input, token)
			if !errors.Is(err, domain.ErrInvalidToken) {
				t.Errorf("expected error %v, got %v", domain.ErrInvalidToken, err)
			}
		})
	}
}

func TestTokenService_CreateAndCacheSecureToken(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	timeGenerator := mocks.NewTimeGeneratorMock(time.Now())
	builder := NewTestBuilder().WithTimeGenerator(timeGenerator).Build()

	tests := map[string]struct {
		input struct {
			tokenType entities.TokenType
			user      *entities.User
		}
		advance     time.Duration
		expectedErr error
	}{
		"get a valid email verification token": {
			input: struct {
				tokenType entities.TokenType
				user      *entities.User
			}{
				tokenType: entities.EmailVerificationToken,
				user:      &entities.User{ID: entities.UserID(uuid.New())},
			},
			advance:     0,
			expectedErr: nil,
		},
		"get an expired email verification token": {
			input: struct {
				tokenType entities.TokenType
				user      *entities.User
			}{
				tokenType: entities.EmailVerificationToken,
				user:      &entities.User{ID: entities.UserID(uuid.New())},
			},
			advance:     emailVerificationTokenExpirationDuration,
			expectedErr: domain.ErrCacheNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Act & Assert
			token, err := builder.TokenService.CreateAndCacheSecureToken(ctx, tt.input.tokenType, tt.input.user)
			if err != nil {
				t.Fatalf("failed to create and cache secure token: %v", err)
			}

			advanceTime(t, builder.TimeGenerator, tt.advance)

			key := utils.GenerateCacheKey(tt.input.tokenType.String(), tt.input.user.ID.String())
			value, err := builder.CacheService.Get(ctx, key)
			if !errors.Is(err, tt.expectedErr) {
				fmt.Println(err, tt.expectedErr)
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil && string(value) != builder.TokenProvider.HashSecureToken(token) {
				t.Errorf("expected token %s, got %s", token, string(value))
			}
		})
	}
}

func TestTokenService_VerifyAndInvalidateSecureToken(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	timeGenerator := mocks.NewTimeGeneratorMock(time.Now())
	builder := NewTestBuilder().WithTimeGenerator(timeGenerator).Build()

	tests := map[string]struct {
		input struct {
			tokenType entities.TokenType
			user      *entities.User
		}
		advance     time.Duration
		expectedErr error
	}{
		"verify and invalidate valid secure token": {
			input: struct {
				tokenType entities.TokenType
				user      *entities.User
			}{
				tokenType: entities.EmailVerificationToken,
				user:      &entities.User{ID: entities.UserID(uuid.New())},
			},
			advance:     0,
			expectedErr: nil,
		},
		"verify and invalidate expired secure token": {
			input: struct {
				tokenType entities.TokenType
				user      *entities.User
			}{
				tokenType: entities.EmailVerificationToken,
				user:      &entities.User{ID: entities.UserID(uuid.New())},
			},
			advance:     emailVerificationTokenExpirationDuration,
			expectedErr: domain.ErrInvalidToken,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Act & Assert
			token, err := builder.TokenService.CreateAndCacheSecureToken(ctx, tt.input.tokenType, tt.input.user)
			if err != nil {
				t.Fatalf("failed to create and cache secure token: %v", err)
			}

			advanceTime(t, builder.TimeGenerator, tt.advance)

			userID, err := builder.TokenService.VerifyAndInvalidateSecureToken(ctx, tt.input.tokenType, token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil && userID != tt.input.user.ID.UUID() {
				t.Errorf("expected user id %s, got %s", tt.input.user.ID.UUID(), userID)
			}

			key := utils.GenerateCacheKey(tt.input.tokenType.String(), tt.input.user.ID.String())

			_, err = builder.CacheService.Get(ctx, key)
			if !errors.Is(err, domain.ErrCacheNotFound) {
				t.Errorf("expected error %v, got %v", domain.ErrCacheNotFound, err)
			}
		})
	}
}
