package services

import (
	"context"
	"errors"
	"go-starter/config"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/utils"
	"sync"
	"time"

	"github.com/google/uuid"
)

// TokenService implements ports.TokenService interface.
// It manages both JWT tokens for authentication and secure tokens for email verification, password reset etc.
type TokenService struct {
	provider          ports.TokenProvider
	cacheSvc          ports.CacheService
	tokenCfg          *config.Token
	tokenTypeDuration *tokenTypeDuration
}

// NewTokenService creates a new instance of TokenService.
func NewTokenService(tokenCfg *config.Token, provider ports.TokenProvider, cacheSvc ports.CacheService) *TokenService {
	return &TokenService{
		provider:          provider,
		cacheSvc:          cacheSvc,
		tokenCfg:          tokenCfg,
		tokenTypeDuration: initTokenTypeDuration(tokenCfg),
	}
}

// tokenTypeDuration represents a thread-safe mapping between token types and their respective durations.
type tokenTypeDuration struct {
	data map[entities.TokenType]time.Duration
	mu   sync.RWMutex
}

// GenerateAccessToken generates a new access token for a user.
// Returns the signed token string or an error if generation fails.
func (ts *TokenService) GenerateAccessToken(user *entities.User) (string, error) {
	token, err := ts.provider.GenerateAccessToken(user, ts.getTokenTypeDuration(entities.AccessToken), ts.tokenCfg.TokenSignature)
	if err != nil {
		return "", domain.ErrInternal
	}
	return token, nil
}

// GenerateRefreshToken generates a new refresh token for a user.
// Returns the signed token string or an error if generation fails.
func (ts *TokenService) GenerateRefreshToken(ctx context.Context, userID entities.UserID) (string, error) {
	tokenID, token, err := ts.provider.GenerateRefreshToken(userID.UUID(), ts.getTokenTypeDuration(entities.RefreshToken), ts.tokenCfg.TokenSignature)
	if err != nil {
		return "", domain.ErrInternal
	}

	key := utils.GenerateCacheKey(entities.RefreshToken.String(), userID.String(), tokenID.String())
	err = ts.cacheSvc.Set(ctx, key, []byte(token), ts.getTokenTypeDuration(entities.RefreshToken))
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyAndParseRefreshToken verifies and parses a refresh token.
// Returns the parsed token claims or an error if validation fails.
func (ts *TokenService) VerifyAndParseRefreshToken(ctx context.Context, token string) (*entities.RefreshTokenClaims, error) {
	claims, err := ts.provider.VerifyAndParseRefreshToken(token, ts.tokenCfg.TokenSignature)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	key := utils.GenerateCacheKey(entities.RefreshToken.String(), claims.Subject.String(), claims.ID.String())
	if _, err := ts.cacheSvc.Get(ctx, key); err != nil {
		if errors.Is(err, domain.ErrCacheNotFound) {
			return nil, domain.ErrInvalidToken
		}
		return nil, err
	}

	return claims, nil
}

func (ts *TokenService) VerifyAndParseAccessToken(token string) (*entities.AccessTokenClaims, error) {
	claims, err := ts.provider.VerifyAndParseAccessToken(token, ts.tokenCfg.TokenSignature)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return claims, nil
}

// RevokeRefreshToken revokes a refresh token by deleting it from the cache.
// Returns an error if the token is not found or if the cache deletion fails.
func (ts *TokenService) RevokeRefreshToken(ctx context.Context, token string) error {
	claims, err := ts.VerifyAndParseRefreshToken(ctx, token)
	if err != nil {
		return err
	}
	key := utils.GenerateCacheKey(entities.RefreshToken.String(), claims.Subject.String(), claims.ID.String())
	return ts.cacheSvc.Delete(ctx, key)
}

// GenerateOneTimeToken generates a new one-time token for a user.
// Returns the token string or an error if generation fails.
func (ts *TokenService) GenerateOneTimeToken(ctx context.Context, tokenType entities.TokenType, userID entities.UserID) (string, error) {
	token, hash, err := ts.provider.GenerateOneTimeToken(userID.UUID())
	if err != nil {
		return "", domain.ErrInternal
	}

	key := utils.GenerateCacheKey(tokenType.String(), userID.String())
	err = ts.cacheSvc.Set(ctx, key, []byte(hash), ts.getTokenTypeDuration(tokenType))
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyAndConsumeOneTimeToken verifies and consumes a one-time token.
// Returns the user ID or an error if the token is not found or if the token is invalid.
func (ts *TokenService) VerifyAndConsumeOneTimeToken(ctx context.Context, tokenType entities.TokenType, token string) (entities.UserID, error) {
	nilUserID := entities.UserID(uuid.Nil)
	parsedToken, err := ts.provider.ParseOneTimeToken(token)
	if err != nil {
		return nilUserID, domain.ErrInvalidToken
	}

	key := utils.GenerateCacheKey(tokenType.String(), parsedToken.UserID.String())
	dbHashedToken, err := ts.cacheSvc.Get(ctx, key)
	if err != nil && errors.Is(err, domain.ErrCacheNotFound) {
		return nilUserID, domain.ErrInvalidToken
	} else if err != nil {
		return nilUserID, err
	}

	if ts.provider.HashOneTimeToken(token) != string(dbHashedToken) {
		return nilUserID, domain.ErrInvalidToken
	}

	err = ts.cacheSvc.Delete(ctx, key)
	if err != nil {
		return nilUserID, err
	}

	return parsedToken.UserID, nil
}

// initTokenTypeDuration initializes a new tokenTypeDuration structure with predefined durations.
func initTokenTypeDuration(tokenCfg *config.Token) *tokenTypeDuration {
	data := map[entities.TokenType]time.Duration{
		entities.AccessToken:            tokenCfg.AccessTokenDuration,
		entities.RefreshToken:           tokenCfg.RefreshTokenDuration,
		entities.EmailVerificationToken: tokenCfg.EmailVerificationTokenDuration,
	}
	return &tokenTypeDuration{
		data: data,
		mu:   sync.RWMutex{},
	}
}

// getTokenTypeDuration returns the duration associated with a specific token type.
func (ts *TokenService) getTokenTypeDuration(tokenType entities.TokenType) time.Duration {
	ts.tokenTypeDuration.mu.RLock()
	defer ts.tokenTypeDuration.mu.RUnlock()
	return ts.tokenTypeDuration.data[tokenType]
}
