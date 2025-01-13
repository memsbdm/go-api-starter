package token

import (
	"github.com/golang-jwt/jwt/v5"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"time"
)

// TokenRepository implements ports.TokenRepository interface and provides access to the JWT package
type TokenRepository struct {
	timeGenerator ports.TimeGenerator
}

// NewTokenRepository creates a new token repository instance
func NewTokenRepository(timeGenerator ports.TimeGenerator) *TokenRepository {
	return &TokenRepository{
		timeGenerator: timeGenerator,
	}
}

// GenerateToken generates a new JWT token for the specified user
func (tr *TokenRepository) GenerateToken(user *entities.User, duration time.Duration, jwtSecret []byte) (string, error) {
	payload := &entities.TokenPayload{
		UserID:    user.ID,
		IssuedAt:  tr.timeGenerator.Now().Unix(),
		ExpiresAt: tr.timeGenerator.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(jwtSecret)
}

// ValidateToken checks if the provided token string is valid and returns the associated claims
func (tr *TokenRepository) ValidateToken(tokenStr string, jwtSecret []byte) (*entities.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &entities.TokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
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

// GenerateRefreshToken creates a new refresh token for a given user
func (tr *TokenRepository) GenerateRefreshToken(user *entities.User, duration time.Duration, jwtSecret []byte) (string, error) {
	token, err := tr.GenerateToken(user, duration, jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

// ValidateRefreshToken validates a refresh token and returns associated token payload
func (tr *TokenRepository) ValidateRefreshToken(tokenStr string, jwtSecret []byte) (*entities.TokenPayload, error) {
	claims, err := tr.ValidateToken(tokenStr, jwtSecret)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
