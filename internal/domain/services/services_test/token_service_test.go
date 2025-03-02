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
	accessTokenExpirationDuration            = 20 * time.Minute
	emailVerificationTokenExpirationDuration = 24 * time.Hour
	passwordResetTokenExpirationDuration     = 15 * time.Minute
)

func TestTokenService_VerifyAuthToken(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		tokenFunc   func(*TestBuilder, *entities.User) (string, error)
		advance     time.Duration
		expectedErr error
	}{
		"verify valid auth token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return builder.TokenService.GenerateAuthToken(context.Background(), user.ID)
			},
			advance:     0,
			expectedErr: nil,
		},
		"verify invalid auth token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return "invalid-token", nil
			},
			advance:     0,
			expectedErr: domain.ErrInvalidToken,
		},
		"verify expired auth token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				return builder.TokenService.GenerateAuthToken(context.Background(), user.ID)
			},
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

			token, err := tt.tokenFunc(builder, user)
			if err != nil {
				t.Fatalf("failed to generate token: %v", err)
			}

			advanceTime(t, builder.TimeGenerator, tt.advance)

			// Act & Assert
			userID, err := builder.TokenService.VerifyAuthToken(context.Background(), token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil && userID != user.ID {
				t.Errorf("expected user id %s, got %s", user.ID, userID)
			}
		})
	}
}

func TestTokenService_RevokeAuthToken(t *testing.T) {
	t.Parallel()

	builder := NewTestBuilder().Build()

	token, err := builder.TokenService.GenerateAuthToken(context.Background(), entities.UserID(uuid.New()))
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	err = builder.TokenService.RevokeAuthToken(context.Background(), token)
	if err != nil {
		t.Fatalf("failed to revoke token: %v", err)
	}

	key := utils.GenerateCacheKey(entities.AccessToken.String(), token)
	_, err = builder.CacheService.Get(context.Background(), key)
	if !errors.Is(err, domain.ErrCacheNotFound) {
		t.Errorf("expected error %v, got %v", domain.ErrCacheNotFound, err)
	}
}

func TestTokenService_GenerateOneTimeToken(t *testing.T) {
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
		"verify and consume overrided token": {
			tokenFunc: func(builder *TestBuilder, user *entities.User) (string, error) {
				overridedToken, err := builder.TokenService.GenerateOneTimeToken(context.Background(), entities.EmailVerificationToken, user.ID)
				if err != nil {
					t.Fatalf("failed to generate token: %v", err)
				}

				_, err = builder.TokenService.GenerateOneTimeToken(context.Background(), entities.EmailVerificationToken, user.ID)
				if err != nil {
					t.Fatalf("failed to generate token: %v", err)
				}

				return overridedToken, nil
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
