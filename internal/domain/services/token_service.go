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

// GenerateJWT generates a new token for the given user.
// Returns the generated token or an error if the generation fails.
func (ts *TokenService) GenerateJWT(tokenType entities.TokenType, user *entities.User) (string, error) {
	token, err := ts.provider.GenerateJWT(tokenType, user, ts.getTokenTypeDuration(tokenType), ts.tokenCfg.TokenSignature)
	if err != nil {
		return "", domain.ErrInternal
	}
	return token, nil
}

// CreateAndCacheJWT creates a new token for the given user and stores it in cache.
// Returns the generated token or an error if the operation fails.
func (ts *TokenService) CreateAndCacheJWT(ctx context.Context, tokenType entities.TokenType, user *entities.User) (string, error) {
	token, err := ts.GenerateJWT(tokenType, user)
	if err != nil {
		return "", domain.ErrInternal
	}

	parsedToken, err := ts.ValidateJWT(tokenType, token)
	if err != nil {
		return "", domain.ErrInvalidToken
	}

	err = ts.cacheToken(ctx, tokenType, user.ID, parsedToken.ID, token)
	if err != nil {
		return "", domain.ErrInternal
	}

	return token, nil
}

// ValidateJWT validates the given token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails.
func (ts *TokenService) ValidateJWT(tokenType entities.TokenType, token string) (*entities.TokenClaims, error) {
	claims, err := ts.provider.ValidateAndParseJWT(tokenType, token, ts.tokenCfg.TokenSignature)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	return claims, nil
}

// VerifyCachedJWT validates the given token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails or if
// the token is not stored in cache.
func (ts *TokenService) VerifyCachedJWT(ctx context.Context, tokenType entities.TokenType, token string) (*entities.TokenClaims, error) {
	claims, err := ts.provider.ValidateAndParseJWT(tokenType, token, ts.tokenCfg.TokenSignature)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	ok, err := ts.isTokenValid(ctx, tokenType, claims.ID, claims.Subject)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

// RevokeJWT deletes the given token from cache.
// Returns an error if the revocation process fails (e.g., if the token is invalid).
func (ts *TokenService) RevokeJWT(ctx context.Context, tokenType entities.TokenType, token string) error {
	claims, err := ts.VerifyCachedJWT(ctx, tokenType, token)
	if err != nil {
		return domain.ErrInvalidToken
	}
	key := generateTokenCacheKey(tokenType, claims.Subject, claims.ID)
	return ts.cacheSvc.Delete(ctx, key)
}

// CreateAndCacheSecureToken creates a new token for the given user and stores it in cache.
// Returns the generated token or an error if the operation fails.
func (ts *TokenService) CreateAndCacheSecureToken(ctx context.Context, tokenType entities.TokenType, user *entities.User) (string, error) {
	token, hashedToken, err := ts.provider.GenerateSecureToken(user.ID.UUID())
	if err != nil {
		return "", domain.ErrInvalidToken
	}

	key := utils.GenerateCacheKey(tokenType.String(), user.ID.String())

	err = ts.cacheSvc.Set(ctx, key, []byte(hashedToken), ts.getTokenTypeDuration(tokenType))
	if err != nil {
		return "", domain.ErrInternal
	}

	return token, nil
}

// VerifyAndInvalidateSecureToken validates the given token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails or if
// the token is not stored in cache.
func (ts *TokenService) VerifyAndInvalidateSecureToken(ctx context.Context, tokenType entities.TokenType, token string) (uuid.UUID, error) {
	parsedToken, err := ts.provider.ParseSecureToken(token)
	if err != nil {
		return uuid.Nil, domain.ErrInvalidToken
	}

	key := utils.GenerateCacheKey(tokenType.String(), parsedToken.UserID.String())
	dbHashedToken, err := ts.cacheSvc.Get(ctx, key)
	if err != nil && !errors.Is(err, domain.ErrCacheNotFound) {
		return uuid.Nil, domain.ErrInternal
	} else if errors.Is(err, domain.ErrCacheNotFound) {
		return uuid.Nil, domain.ErrInvalidToken
	}

	userHashedToken := ts.provider.HashSecureToken(token)
	if userHashedToken != string(dbHashedToken) {
		return uuid.Nil, domain.ErrInvalidToken
	}

	err = ts.cacheSvc.DeleteByPrefix(ctx, key)
	if err != nil {
		return uuid.Nil, domain.ErrInternal
	}

	return parsedToken.UserID, nil
}

// cacheToken stores the token in the cache associated with the given user ID and token ID.
// It constructs a unique cache key and sets the token with an expiration duration.
func (ts *TokenService) cacheToken(ctx context.Context, tokenType entities.TokenType, userID entities.UserID, tokenID uuid.UUID, token string) error {
	key := generateTokenCacheKey(tokenType, userID, tokenID)
	return ts.cacheSvc.Set(ctx, key, []byte(token), ts.getTokenTypeDuration(tokenType))
}

// isTokenValid checks if the token is present in the cache for the given user ID and token ID.
// Returns a boolean indicating presence and an error if the operation fails.
func (ts *TokenService) isTokenValid(ctx context.Context, tokenType entities.TokenType, tokenID uuid.UUID, userID entities.UserID) (bool, error) {
	key := generateTokenCacheKey(tokenType, userID, tokenID)
	if _, err := ts.cacheSvc.Get(ctx, key); err != nil {
		if errors.Is(err, domain.ErrCacheNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// generateTokenCacheKey creates a unique cache key by combining the token type, user ID, and token ID.
func generateTokenCacheKey(tokenType entities.TokenType, userID entities.UserID, tokenID uuid.UUID) string {
	return utils.GenerateCacheKey(tokenType.String(), userID.String(), tokenID)
}
