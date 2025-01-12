package token

import (
	"github.com/golang-jwt/jwt/v5"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"time"
)

// JWTTokenImpl implements ports.TokenRepository interface and provides access to the JWT package
type JWTTokenImpl struct {
}

// NewJWTTokenImpl creates a new token repository instance
func NewJWTTokenImpl() *JWTTokenImpl {
	return &JWTTokenImpl{}
}

// GenerateToken generates a new JWT token for the specified user
func (jt *JWTTokenImpl) GenerateToken(user *entities.User, duration time.Duration, jwtSecret []byte) (string, error) {
	payload := &entities.TokenPayload{
		UserID:    user.ID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(jwtSecret)
}

// ValidateToken checks if the provided token string is valid and returns the associated claims
func (jt *JWTTokenImpl) ValidateToken(tokenStr string, jwtSecret []byte) (*entities.TokenPayload, error) {
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
func (jt *JWTTokenImpl) GenerateRefreshToken(user *entities.User, duration time.Duration, jwtSecret []byte) (string, error) {
	token, err := jt.GenerateToken(user, duration, jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateRefreshToken validates a refresh token and returns associated token payload
func (jt *JWTTokenImpl) ValidateRefreshToken(tokenStr string, jwtSecret []byte) (*entities.TokenPayload, error) {
	claims, err := jt.ValidateToken(tokenStr, jwtSecret)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
