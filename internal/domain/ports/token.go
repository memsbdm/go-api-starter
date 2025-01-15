package ports

import (
	"context"
	"go-starter/internal/domain/entities"
	"time"
)

// TokenService is an interface for interacting with token-related business logic
type TokenService interface {
	GenerateAccessToken(user *entities.User) (string, error)
	ValidateAndParseAccessToken(token string) (*entities.AccessTokenClaims, error)
	GenerateRefreshToken(ctx context.Context, userID entities.UserID) (string, error)
	ValidateAndParseRefreshToken(ctx context.Context, token string) (*entities.RefreshTokenClaims, error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
}

type TokenRepository interface {
	GenerateAccessToken(user *entities.User, duration time.Duration, signature []byte) (string, error)
	ValidateAndParseAccessToken(token string, signature []byte) (*entities.AccessTokenClaims, error)
	GenerateRefreshToken(userID entities.UserID, duration time.Duration, signature []byte) (entities.RefreshTokenID, string, error)
	ValidateAndParseRefreshToken(token string, signature []byte) (*entities.RefreshTokenClaims, error)
}
