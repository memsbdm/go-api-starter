package services

import (
	"context"
	"go-starter/config"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/utils"
)

// TokenService implements ports.TokenService interface.
type TokenService struct {
	repo       ports.TokenRepository
	cacheSvc   ports.CacheService
	cfg        *config.Token
	errTracker ports.ErrorTracker
}

// NewTokenService creates a new instance of TokenService.
func NewTokenService(cfg *config.Token, repo ports.TokenRepository, cacheSvc ports.CacheService, errTracker ports.ErrorTracker) *TokenService {
	return &TokenService{
		repo:       repo,
		cacheSvc:   cacheSvc,
		cfg:        cfg,
		errTracker: errTracker,
	}
}

// RefreshTokenCachePrefix is the prefix for caching refresh tokens.
const RefreshTokenCachePrefix = "refresh_token"

// GenerateAccessToken generates a new access token for the given user.
// Returns the generated access token or an error if the generation fails.
func (ts *TokenService) GenerateAccessToken(user *entities.User) (entities.AccessToken, error) {
	token, err := ts.repo.GenerateAccessToken(user, ts.cfg.AccessTokenDuration, ts.cfg.AccessTokenSignature)
	if err != nil {
		return "", domain.ErrInternal
	}
	return token, nil
}

// ValidateAndParseAccessToken validates the given access token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails.
func (ts *TokenService) ValidateAndParseAccessToken(token string) (*entities.AccessTokenClaims, error) {
	tokenPayload, err := ts.repo.ValidateAndParseAccessToken(token, ts.cfg.AccessTokenSignature)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	return tokenPayload, nil
}

// GenerateRefreshToken creates a new refresh token for the given user ID.
// Returns the generated refresh token or an error if the operation fails.
func (ts *TokenService) GenerateRefreshToken(ctx context.Context, userID entities.UserID) (entities.RefreshToken, error) {
	tokenID, token, err := ts.repo.GenerateRefreshToken(userID, ts.cfg.AccessTokenDuration, ts.cfg.RefreshTokenSignature)
	if err != nil {
		return "", domain.ErrInternal
	}

	err = ts.storeRefreshTokenInCache(ctx, userID, tokenID, token)
	if err != nil {
		return "", domain.ErrInternal
	}

	return token, nil
}

// ValidateAndParseRefreshToken validates the given refresh token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails.
func (ts *TokenService) ValidateAndParseRefreshToken(ctx context.Context, token string) (*entities.RefreshTokenClaims, error) {
	claims, err := ts.repo.ValidateAndParseRefreshToken(token, ts.cfg.RefreshTokenSignature)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	ok, err := ts.isRefreshTokenInCache(ctx, claims.ID, claims.Subject)
	if err != nil || !ok {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

// RevokeRefreshToken invalidates the given refresh token.
// Returns an error if the revocation process fails (e.g., if the token is invalid).
func (ts *TokenService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	claims, err := ts.ValidateAndParseRefreshToken(ctx, refreshToken)
	if err != nil {
		return domain.ErrInvalidToken
	}
	key := utils.GenerateCacheKey(RefreshTokenCachePrefix, claims.Subject, claims.ID)
	return ts.cacheSvc.Delete(ctx, key)
}

// storeRefreshTokenInCache stores the refresh token in the cache associated with the given user ID and refresh token ID.
// It constructs a unique cache key and sets the refresh token with an expiration duration.
func (ts *TokenService) storeRefreshTokenInCache(ctx context.Context, userID entities.UserID, refreshTokenID entities.RefreshTokenID, refreshToken entities.RefreshToken) error {
	key := utils.GenerateCacheKey(RefreshTokenCachePrefix, userID.String(), refreshTokenID.String())
	return ts.cacheSvc.Set(ctx, key, []byte(refreshToken), ts.cfg.RefreshTokenDuration)
}

// isRefreshTokenInCache checks if the refresh token is present in the cache for the given user ID and refresh token ID.
// Returns a boolean indicating presence and an error if the operation fails.
func (ts *TokenService) isRefreshTokenInCache(ctx context.Context, refreshTokenID entities.RefreshTokenID, userID entities.UserID) (bool, error) {
	key := utils.GenerateCacheKey(RefreshTokenCachePrefix, userID.String(), refreshTokenID.String())
	if value, err := ts.cacheSvc.Get(ctx, key); err != nil {
		return false, domain.ErrInternal
	} else if value == nil {
		return false, nil
	}

	return true, nil
}
