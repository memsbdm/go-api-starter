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

// GenerateAccessToken generates a new access token for the specified user
func (ts *TokenService) GenerateAccessToken(user *entities.User) (string, error) {
	token, err := ts.repo.GenerateToken(user, ts.cfg.AccessTokenDuration, ts.cfg.JWTSecret)
	if err != nil {
		return "", domain.ErrInternal
	}
	return token, nil
}

// ValidateAccessToken checks if the provided access token is valid and returns the associated claims
func (ts *TokenService) ValidateAccessToken(tokenStr string) (*entities.TokenPayload, error) {
	tokenPayload, err := ts.repo.ValidateToken(tokenStr, ts.cfg.JWTSecret)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	return tokenPayload, nil
}

// GenerateRefreshToken creates a new refresh token for a given user and stores it in cache
func (ts *TokenService) GenerateRefreshToken(user *entities.User) (string, error) {
	token, err := ts.repo.GenerateRefreshToken(user, ts.cfg.RefreshTokenDuration, ts.cfg.JWTSecret)
	if err != nil {
		return "", domain.ErrInternal
	}

	ctx := context.Background()
	key := utils.GenerateCacheKey("refresh_token", user.ID)

	// Delete previous refresh token if exists
	err = ts.cacheSvc.Delete(ctx, key)
	if err != nil {
		return "", domain.ErrInternal
	}

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
	claims, err := ts.repo.ValidateRefreshToken(refreshToken, ts.cfg.JWTSecret)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	ctx := context.Background()
	key := utils.GenerateCacheKey("refresh_token", claims.UserID)
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
