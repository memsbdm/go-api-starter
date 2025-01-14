package ports

import (
	"go-starter/internal/domain/entities"
	"time"
)

// TokenService is an interface for interacting with token-related business logic
type TokenService interface {
	// GenerateAccessToken creates a new access token for a given user
	GenerateAccessToken(user *entities.User) (string, error)
	// ValidateAccessToken validates an access token and returns associated token payload
	ValidateAccessToken(tokenStr string) (*entities.TokenPayload, error)
	// GenerateRefreshToken creates a new refresh token for a given user
	GenerateRefreshToken(user *entities.User) (string, error)
	// ValidateRefreshToken validates a refresh token and returns associated token payload
	ValidateRefreshToken(tokenStr string) (*entities.TokenPayload, error)
}

type TokenRepository interface {
	// GenerateToken creates a new token for a given user
	GenerateToken(user *entities.User, duration time.Duration, jwtSecret []byte) (string, error)
	// ValidateToken validates an token token and returns associated token payload
	ValidateToken(tokenStr string, jwtSecret []byte) (*entities.TokenPayload, error)
	// GenerateRefreshToken creates a new refresh token for a given user
	GenerateRefreshToken(user *entities.User, duration time.Duration, jwtSecret []byte) (string, error)
	// ValidateRefreshToken validates a refresh token and returns associated token payload
	ValidateRefreshToken(tokenStr string, jwtSecret []byte) (*entities.TokenPayload, error)
}
