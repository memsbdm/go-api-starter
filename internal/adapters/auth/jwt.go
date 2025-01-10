package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"go-starter/config"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"time"
)

// JWTToken implements ports.TokenService interface
type JWTToken struct {
	jwtSecret []byte
	duration  time.Duration
}

// NewTokenService creates a new instance of JWTToken based on the provided configuration
func NewTokenService(cfg *config.Token) ports.TokenService {
	return &JWTToken{
		jwtSecret: cfg.JWTSecret,
		duration:  cfg.Duration,
	}
}

// CreateToken generates a new JWT token for the specified user
func (jt *JWTToken) CreateToken(user *entities.User) (string, error) {
	payload := &entities.TokenPayload{
		UserID:    user.ID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(jt.duration).Unix(),
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
