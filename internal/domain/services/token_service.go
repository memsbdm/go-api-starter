package services

import (
	"context"
	"github.com/google/uuid"
	"go-starter/config"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/utils"
	"sync"
	"time"
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
		entities.AccessToken:  tokenCfg.AccessTokenDuration,
		entities.RefreshToken: tokenCfg.RefreshTokenDuration,
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

// Generate generates a new token for the given user.
// Returns the generated token or an error if the generation fails.
func (ts *TokenService) Generate(tokenType entities.TokenType, user *entities.User) (string, error) {
	token, err := ts.provider.Generate(tokenType, user, ts.getTokenTypeDuration(tokenType), ts.tokenCfg.TokenSignature)
	if err != nil {
		return "", domain.ErrInternal
	}
	return token, nil
}

// GenerateTokenWithCache creates a new token for the given user and stores it in cache.
// Returns the generated token or an error if the operation fails.
func (ts *TokenService) GenerateTokenWithCache(ctx context.Context, tokenType entities.TokenType, user *entities.User) (string, error) {
	token, err := ts.Generate(tokenType, user)
	if err != nil {
		return "", domain.ErrInternal
	}

	parsedToken, _ := ts.ValidateAndParse(tokenType, token)

	err = ts.storeInCache(ctx, tokenType, user.ID, parsedToken.ID, token)
	if err != nil {
		return "", domain.ErrInternal
	}

	return token, nil
}

// ValidateAndParse validates the given token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails.
func (ts *TokenService) ValidateAndParse(tokenType entities.TokenType, token string) (*entities.TokenClaims, error) {
	claims, err := ts.provider.ValidateAndParse(tokenType, token, ts.tokenCfg.TokenSignature)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	return claims, nil
}

// ValidateAndParseWithCache validates the given token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails or if
// the token is not stored in cache.
func (ts *TokenService) ValidateAndParseWithCache(ctx context.Context, tokenType entities.TokenType, token string) (*entities.TokenClaims, error) {
	claims, err := ts.provider.ValidateAndParse(tokenType, token, ts.tokenCfg.TokenSignature)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	ok, err := ts.isTokenInCache(ctx, tokenType, claims.ID, claims.Subject)
	if err != nil || !ok {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

// RevokeTokenFromCache deletes the given token from cache.
// Returns an error if the revocation process fails (e.g., if the token is invalid).
func (ts *TokenService) RevokeTokenFromCache(ctx context.Context, tokenType entities.TokenType, token string) error {
	claims, err := ts.ValidateAndParseWithCache(ctx, tokenType, token)
	if err != nil {
		return domain.ErrInvalidToken
	}
	key := generateTokenCacheKey(tokenType, claims.Subject, claims.ID)
	return ts.cacheSvc.Delete(ctx, key)
}

// storeRefreshTokenInCache stores the refresh token in the cache associated with the given user ID and refresh token ID.
// It constructs a unique cache key and sets the refresh token with an expiration duration.
func (ts *TokenService) storeInCache(ctx context.Context, tokenType entities.TokenType, userID entities.UserID, tokenID uuid.UUID, token string) error {
	key := generateTokenCacheKey(tokenType, userID, tokenID)
	return ts.cacheSvc.Set(ctx, key, []byte(token), ts.getTokenTypeDuration(tokenType))
}

// isTokenInCache checks if the token is present in the cache for the given user ID and token ID.
// Returns a boolean indicating presence and an error if the operation fails.
func (ts *TokenService) isTokenInCache(ctx context.Context, tokenType entities.TokenType, tokenID uuid.UUID, userID entities.UserID) (bool, error) {
	key := generateTokenCacheKey(tokenType, userID, tokenID)
	if value, err := ts.cacheSvc.Get(ctx, key); err != nil {
		return false, domain.ErrInternal
	} else if value == nil {
		return false, nil
	}

	return true, nil
}

// generateTokenCacheKey creates a unique cache key by combining the token type, user ID, and token ID.
func generateTokenCacheKey(tokenType entities.TokenType, userID entities.UserID, tokenID uuid.UUID) string {
	return utils.GenerateCacheKey(tokenType.String(), userID.String(), tokenID)
}
