package services

import (
	"context"
	"go-starter/config"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/utils"
)

// TokenService implements ports.TokenService interface and provides access to the token repository
type TokenService struct {
	repo     ports.TokenRepository
	cacheSvc ports.CacheService
	cfg      *config.Token
}

// NewTokenService creates a new token service instance
func NewTokenService(cfg *config.Token, repo ports.TokenRepository, cacheSvc ports.CacheService) *TokenService {
	return &TokenService{
		repo:     repo,
		cacheSvc: cacheSvc,
		cfg:      cfg,
	}
}

const RefreshTokenCachePrefix = "refresh_token"

// GenerateAccessToken generates a new access token for the specified user
func (ts *TokenService) GenerateAccessToken(user *entities.User) (string, error) {
	token, err := ts.repo.GenerateAccessToken(user, ts.cfg.AccessTokenDuration, ts.cfg.AccessTokenSignature)
	if err != nil {
		return "", domain.ErrInternal
	}
	return token, nil
}

func (ts *TokenService) ValidateAndParseAccessToken(token string) (*entities.AccessTokenClaims, error) {
	tokenPayload, err := ts.repo.ValidateAndParseAccessToken(token, ts.cfg.AccessTokenSignature)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	return tokenPayload, nil
}

func (ts *TokenService) GenerateRefreshToken(ctx context.Context, userID entities.UserID) (string, error) {
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

func (ts *TokenService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	claims, err := ts.ValidateAndParseRefreshToken(ctx, refreshToken)
	if err != nil {
		return domain.ErrInvalidToken
	}
	key := utils.GenerateCacheKey(RefreshTokenCachePrefix, claims.Subject, claims.ID)
	return ts.cacheSvc.Delete(ctx, key)
}

func (ts *TokenService) storeRefreshTokenInCache(ctx context.Context, userID entities.UserID, refreshTokenID entities.RefreshTokenID, refreshToken string) error {
	key := utils.GenerateCacheKey(RefreshTokenCachePrefix, userID.String(), refreshTokenID.String())
	return ts.cacheSvc.Set(ctx, key, []byte(refreshToken), ts.cfg.RefreshTokenDuration)
}

func (ts *TokenService) isRefreshTokenInCache(ctx context.Context, refreshTokenID entities.RefreshTokenID, userID entities.UserID) (bool, error) {
	key := utils.GenerateCacheKey(RefreshTokenCachePrefix, userID.String(), refreshTokenID.String())
	if value, err := ts.cacheSvc.Get(ctx, key); err != nil {
		return false, domain.ErrInternal
	} else if value == nil {
		return false, nil
	}

	return true, nil
}
