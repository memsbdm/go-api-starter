package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"go-starter/config"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/utils"
	"time"
)

// JWTToken implements ports.TokenService interface
type JWTToken struct {
	jwtSecret            []byte
	tokenDuration        time.Duration
	refreshTokenDuration time.Duration
	cacheSvc             ports.CacheService
}

// NewTokenService creates a new instance of JWTToken based on the provided configuration
func NewTokenService(cfg *config.Token, cacheSvc ports.CacheService) ports.TokenService {
	return &JWTToken{
		jwtSecret:            cfg.JWTSecret,
		tokenDuration:        cfg.TokenDuration,
		refreshTokenDuration: cfg.RefreshTokenDuration,
		cacheSvc:             cacheSvc,
	}
}

// CreateToken generates a new JWT token for the specified user
func (jt *JWTToken) CreateToken(user *entities.User, isRefreshToken bool) (string, error) {
	duration := jt.tokenDuration
	if isRefreshToken {
		duration = jt.refreshTokenDuration
	}

	payload := &entities.TokenPayload{
		UserID:    user.ID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(jt.jwtSecret)
}

// ValidateToken checks if the provided token string is valid and returns the associated claims
func (jt *JWTToken) ValidateToken(tokenStr string) (*entities.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &entities.TokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		return jt.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	tokenPayload, ok := token.Claims.(*entities.TokenPayload)
	if !ok || !token.Valid {
		return nil, domain.ErrUnauthorized
	}

	return tokenPayload, nil
}

// CreateRefreshToken creates a new refresh token for a given user
func (jt *JWTToken) CreateRefreshToken(user *entities.User) (string, error) {
	token, err := jt.CreateToken(user, true)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	key := utils.GenerateCacheKey("refresh_token", user.ID)

	// Delete previous refresh token if exists
	err = jt.cacheSvc.Delete(ctx, key)
	if err != nil {
		return "", err
	}

	value, err := utils.Serialize(token)
	if err != nil {
		return "", err
	}

	err = jt.cacheSvc.Set(ctx, key, value, jt.refreshTokenDuration)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateRefreshToken validates a refresh token and returns associated token payload
func (jt *JWTToken) ValidateRefreshToken(tokenStr string) error {
	claims, err := jt.ValidateToken(tokenStr)
	if err != nil {
		return err
	}

	ctx := context.Background()
	key := utils.GenerateCacheKey("refresh_token", claims.UserID)

	_, err = jt.cacheSvc.Get(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
