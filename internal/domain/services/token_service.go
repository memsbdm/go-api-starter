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
	_, token, err := ts.repo.GenerateToken(user, ts.cfg.AccessTokenDuration, ts.cfg.JWTSecret)
	if err != nil {
		return "", domain.ErrInternal
	}
	return token, nil
}

// GetTokenPayload checks if the provided token is valid and returns the associated token payload
func (ts *TokenService) GetTokenPayload(tokenStr string) (*entities.TokenPayload, error) {
	tokenPayload, err := ts.repo.ValidateToken(tokenStr, ts.cfg.JWTSecret)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	return tokenPayload, nil
}

// GenerateRefreshToken creates a new refresh token for a given user and stores it in cache
func (ts *TokenService) GenerateRefreshToken(user *entities.User) (string, error) {
	ctx := context.Background()

	// Generate a new token
	tokenID, token, err := ts.repo.GenerateRefreshToken(user, ts.cfg.RefreshTokenDuration, ts.cfg.JWTSecret)
	if err != nil {
		return "", domain.ErrInternal
	}

	// Store in cache
	key := utils.GenerateCacheKey(RefreshTokenCachePrefix, user.ID, tokenID)
	value, err := utils.Serialize(token)
	if err != nil {
		return "", domain.ErrInternal
	}

	err = ts.cacheSvc.Set(ctx, key, value, ts.cfg.RefreshTokenDuration)
	if err != nil {
		return "", domain.ErrInternal
	}

	return token, nil
}

// ValidateRefreshToken validates a refresh token and returns associated token payload
func (ts *TokenService) ValidateRefreshToken(refreshToken string) (*entities.TokenPayload, error) {
	ctx := context.Background()

	claims, err := ts.repo.ValidateRefreshToken(refreshToken, ts.cfg.JWTSecret)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Verify token with cached one
	key := utils.GenerateCacheKey(RefreshTokenCachePrefix, claims.UserID, claims.ID)
	value, err := ts.cacheSvc.Get(ctx, key)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	var storedRefreshToken string
	err = utils.Deserialize(value, &storedRefreshToken)
	if err != nil {
		return nil, domain.ErrInternal
	}

	if refreshToken != storedRefreshToken {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

func (ts *TokenService) DeleteRefreshToken(ctx context.Context, tokenStr string) error {
	tokenPayload, err := ts.GetTokenPayload(tokenStr)
	if err != nil {
		return domain.ErrInvalidToken
	}
	key := utils.GenerateCacheKey(RefreshTokenCachePrefix, tokenPayload.UserID, tokenPayload.ID)
	return ts.cacheSvc.Delete(ctx, key)
}
