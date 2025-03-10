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
)

// TokenService implements ports.TokenService interface.
// It manages one-time and authentication tokens.
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

// GenerateAuthToken generates an access token for a user.
// Returns the access token or an error if the operation fails.
func (ts *TokenService) GenerateAuthToken(ctx context.Context, userID entities.UserID) (string, error) {
	token, err := ts.provider.GenerateRandomToken()
	if err != nil {
		return "", domain.ErrInternal
	}

	key := utils.GenerateCacheKey(entities.AccessToken.String(), token)

	err = ts.cacheSvc.Set(ctx, key, []byte(userID.String()), ts.getTokenTypeDuration(entities.AccessToken))
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyAuthToken verifies an access token.
// Returns the user ID or an error if the token is not found or if the token is invalid.
func (ts *TokenService) VerifyAuthToken(ctx context.Context, token string) (entities.UserID, error) {
	key := utils.GenerateCacheKey(entities.AccessToken.String(), token)
	userIDBytes, err := ts.cacheSvc.Get(ctx, key)
	if err != nil {
		if errors.Is(err, domain.ErrCacheNotFound) {
			return entities.NilUserID, domain.ErrInvalidToken
		}
		return entities.NilUserID, err
	}

	err = ts.cacheSvc.Set(ctx, key, userIDBytes, ts.getTokenTypeDuration(entities.AccessToken))
	if err != nil {
		return entities.NilUserID, err
	}

	userID, err := entities.ParseUserID(string(userIDBytes))
	if err != nil {
		return entities.NilUserID, domain.ErrInternal
	}

	return userID, nil
}

// RevokeAuthToken revokes an access token.
// Returns an error if the revocation fails.
func (ts *TokenService) RevokeAuthToken(ctx context.Context, token string) error {
	key := utils.GenerateCacheKey(entities.AccessToken.String(), token)
	return ts.cacheSvc.Delete(ctx, key)
}

// GenerateOneTimeToken generates a new one-time token for a user.
// Returns the token string or an error if generation fails.
func (ts *TokenService) GenerateOneTimeToken(ctx context.Context, tokenType entities.TokenType, userID entities.UserID) (string, error) {
	token, err := ts.provider.GenerateOneTimeToken(userID)
	if err != nil {
		return "", domain.ErrInternal
	}

	key := utils.GenerateCacheKey(tokenType.String(), userID.String())
	err = ts.cacheSvc.Set(ctx, key, []byte(token), ts.getTokenTypeDuration(tokenType))
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyAndConsumeOneTimeToken verifies and consumes a one-time token.
// Returns the user ID or an error if the token is not found in cache or if the token is invalid.
func (ts *TokenService) VerifyAndConsumeOneTimeToken(ctx context.Context, tokenType entities.TokenType, token string) (entities.UserID, error) {
	userID, err := ts.VerifyOneTimeToken(ctx, tokenType, token)
	if err != nil {
		return entities.NilUserID, err
	}

	return userID, ts.ConsumeOneTimeToken(ctx, tokenType, token)
}

// VerifyOneTimeToken verifies a one-time token.
// Returns the user ID or an error if the token is not found in cache or if the token is invalid.
func (ts *TokenService) VerifyOneTimeToken(ctx context.Context, tokenType entities.TokenType, token string) (entities.UserID, error) {
	userID, err := ts.provider.ParseOneTimeToken(token)
	if err != nil {
		return entities.NilUserID, domain.ErrInvalidToken
	}

	key := utils.GenerateCacheKey(tokenType.String(), userID.String())
	dbToken, err := ts.cacheSvc.Get(ctx, key)
	if err != nil && errors.Is(err, domain.ErrCacheNotFound) {
		return entities.NilUserID, domain.ErrInvalidToken
	} else if err != nil {
		return entities.NilUserID, err
	}

	if string(dbToken) != token {
		return entities.NilUserID, domain.ErrInvalidToken
	}

	return userID, nil
}

// ConsumeOneTimeToken consumes a one-time token.
// Returns an error if the consumption fails.
func (ts *TokenService) ConsumeOneTimeToken(ctx context.Context, tokenType entities.TokenType, token string) error {
	userID, err := ts.VerifyOneTimeToken(ctx, tokenType, token)
	if err != nil {
		return err
	}

	key := utils.GenerateCacheKey(tokenType.String(), userID.String())
	return ts.cacheSvc.Delete(ctx, key)
}

// initTokenTypeDuration initializes a new tokenTypeDuration structure with predefined durations.
func initTokenTypeDuration(tokenCfg *config.Token) *tokenTypeDuration {
	data := map[entities.TokenType]time.Duration{
		entities.AccessToken:            tokenCfg.AccessTokenDuration,
		entities.EmailVerificationToken: tokenCfg.EmailVerificationTokenDuration,
		entities.PasswordResetToken:     tokenCfg.PasswordResetTokenDuration,
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
