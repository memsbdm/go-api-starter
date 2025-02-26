//go:build !integration

package services_test

import (
	"context"
	"errors"
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

func TestTokenService_VerifyAndParseAccessToken(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		tokenFunc   func(*TestBuilder, *entities.User) (string, error)
		advance     time.Duration
		expectedErr error
	}{
		"verify and parse valid access token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return builder.TokenService.GenerateAccessToken(user)
			},
			advance:     0,
			expectedErr: nil,
		},
		"verify and parse expired access token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return builder.TokenService.GenerateAccessToken(user)
			},
			advance:     accessTokenExpirationDuration,
			expectedErr: domain.ErrInvalidToken,
		},
		"verify and parse invalid token format": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return "invalid-token-format", nil
			},
			advance:     0,
			expectedErr: domain.ErrInvalidToken,
		},
		"verify and parse refresh token instead of access token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return builder.TokenService.GenerateRefreshToken(context.Background(), user.ID)
			},
			advance:     0,
			expectedErr: domain.ErrInvalidToken,
		},
		"verify and parse tampered token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				token, err := builder.TokenService.GenerateAccessToken(user)
				if err != nil {
					return "", err
				}
				if len(token) > 0 {
					return token[:len(token)-1] + "X", nil
				}
				return token, nil
			},
			advance:     0,
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

			token, err := tt.tokenFunc(builder, user)
			if err != nil {
				t.Fatalf("failed to generate token: %v", err)
			}

			advanceTime(t, builder.TimeGenerator, tt.advance)

			// Act & Assert
			claims, err := builder.TokenService.VerifyAndParseAccessToken(token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil && claims.Subject != user.ID {
				t.Errorf("expected user id %s, got %s", user.ID, claims.Subject)
			}
		})
	}
}

func TestTokenService_VerifyAndParseRefreshToken(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		tokenFunc   func(*TestBuilder, *entities.User) (string, error)
		advance     time.Duration
		expectedErr error
	}{
		"verify and parse valid refresh token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return builder.TokenService.GenerateRefreshToken(context.Background(), user.ID)
			},
			advance:     0,
			expectedErr: nil,
		},
		"verify and parse expired refresh token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return builder.TokenService.GenerateRefreshToken(context.Background(), user.ID)
			},
			advance:     refreshTokenExpirationDuration,
			expectedErr: domain.ErrInvalidToken,
		},
		"verify and parse invalid token format": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return "invalid-token-format", nil
			},
			advance:     0,
			expectedErr: domain.ErrInvalidToken,
		},
		"verify and parse access token instead of refresh token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return builder.TokenService.GenerateAccessToken(user)
			},
			advance:     0,
			expectedErr: domain.ErrInvalidToken,
		},
		"verify and parse tampered token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				token, err := builder.TokenService.GenerateRefreshToken(context.Background(), user.ID)
				if err != nil {
					return "", err
				}
				if len(token) > 0 {
					return token[:len(token)-1] + "X", nil
				}
				return token, nil
			},
			advance:     0,
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

			token, err := tt.tokenFunc(builder, user)
			if err != nil {
				t.Fatalf("failed to generate token: %v", err)
			}

			advanceTime(t, builder.TimeGenerator, tt.advance)

			// Act & Assert
			claims, err := builder.TokenService.VerifyAndParseRefreshToken(context.Background(), token)
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

	tests := map[string]struct {
		tokenFunc   func(*TestBuilder, *entities.User) (string, error)
		expectedErr error
	}{
		"revoke valid refresh token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return builder.TokenService.GenerateRefreshToken(context.Background(), user.ID)
			},
			expectedErr: nil,
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

			token, err := tt.tokenFunc(builder, user)
			if err != nil {
				t.Fatalf("failed to generate token: %v", err)
			}

			// Act & Assert
			err = builder.TokenService.RevokeRefreshToken(context.Background(), token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			key := utils.GenerateCacheKey(entities.RefreshToken.String(), user.ID.String(), token)
			_, err = builder.CacheService.Get(context.Background(), key)
			if !errors.Is(err, domain.ErrCacheNotFound) {
				t.Errorf("expected error %v, got %v", domain.ErrCacheNotFound, err)
			}
		})
	}
}

func TestTokenService_VerifyAndConsumeOneTimeToken(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		tokenFunc   func(*TestBuilder, *entities.User) (string, error)
		advance     time.Duration
		expectedErr error
	}{
		"verify and consume valid one-time token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return builder.TokenService.GenerateOneTimeToken(context.Background(), entities.EmailVerificationToken, user.ID)
			},
			advance:     0,
			expectedErr: nil,
		},
		"verify and consume expired one-time token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return builder.TokenService.GenerateOneTimeToken(context.Background(), entities.EmailVerificationToken, user.ID)
			},
			advance:     emailVerificationTokenExpirationDuration,
			expectedErr: domain.ErrInvalidToken,
		},
		"verify and consume invalid token format": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return "invalid-token-format", nil
			},
			advance:     0,
			expectedErr: domain.ErrInvalidToken,
		},
		"verify and consume tampered token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				token, err := builder.TokenService.GenerateOneTimeToken(context.Background(), entities.EmailVerificationToken, user.ID)
				if err != nil {
					return "", err
				}
				if len(token) > 0 {
					return token[:len(token)-1] + "X", nil
				}
				return token, nil
			},
			advance:     0,
			expectedErr: domain.ErrInvalidToken,
		},
		"verify and consume invalid token type": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return builder.TokenService.GenerateOneTimeToken(context.Background(), "invalid-token-type", user.ID)
			},
			advance:     0,
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

			token, err := tt.tokenFunc(builder, user)
			if err != nil {
				t.Fatalf("failed to generate token: %v", err)
			}

			advanceTime(t, builder.TimeGenerator, tt.advance)

			// Act & Assert
			userID, err := builder.TokenService.VerifyAndConsumeOneTimeToken(context.Background(), entities.EmailVerificationToken, token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil && userID != user.ID {
				t.Errorf("expected user id %s, got %s", user.ID, userID)
			}

			key := utils.GenerateCacheKey(entities.EmailVerificationToken.String(), user.ID.String(), token)
			_, err = builder.CacheService.Get(context.Background(), key)
			if !errors.Is(err, domain.ErrCacheNotFound) {
				t.Errorf("expected error %v, got %v", domain.ErrCacheNotFound, err)
			}
		})
	}
}
