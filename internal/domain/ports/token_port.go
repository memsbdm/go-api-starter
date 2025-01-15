package ports

import (
	"context"
	"go-starter/internal/domain/entities"
	"time"
)

// TokenService is an interface for interacting with token-related business logic.
type TokenService interface {
	// GenerateAccessToken generates a new access token for the given user.
	// Returns the generated access token as a string or an error if the generation fails.
	GenerateAccessToken(user *entities.User) (string, error)

	// ValidateAndParseAccessToken validates the given access token and extracts its claims.
	// Returns a structured representation of the token claims or an error if validation fails.
	ValidateAndParseAccessToken(token string) (*entities.AccessTokenClaims, error)

	// GenerateRefreshToken creates a new refresh token for the given user ID.
	// Returns the generated refresh token as a string or an error if the operation fails.
	GenerateRefreshToken(ctx context.Context, userID entities.UserID) (string, error)

	// ValidateAndParseRefreshToken validates the given refresh token and extracts its claims.
	// Returns a structured representation of the token claims or an error if validation fails.
	ValidateAndParseRefreshToken(ctx context.Context, token string) (*entities.RefreshTokenClaims, error)

	// RevokeRefreshToken invalidates the given refresh token.
	// Returns an error if the revocation process fails (e.g., if the token is invalid).
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
}

// TokenRepository is an interface for interacting with token-related data and cryptographic operations.
type TokenRepository interface {
	// GenerateAccessToken generates a new JWT access token for the given user.
	// Returns the generated access token as a string or an error if the generation fails.
	GenerateAccessToken(user *entities.User, duration time.Duration, signature []byte) (string, error)

	// ValidateAndParseAccessToken validates the given JWT access token and extracts its claims.
	// Returns a structured representation of the token claims or an error if validation fails.
	ValidateAndParseAccessToken(token string, signature []byte) (*entities.AccessTokenClaims, error)

	// GenerateRefreshToken creates a new JWT refresh token for the given user ID.
	// Returns a unique refresh token ID, the token string, or an error if the operation fails.
	GenerateRefreshToken(userID entities.UserID, duration time.Duration, signature []byte) (entities.RefreshTokenID, string, error)

	// ValidateAndParseRefreshToken validates the given JWT refresh token and extracts its claims.
	// Returns a structured representation of the token claims or an error if validation fails.
	ValidateAndParseRefreshToken(token string, signature []byte) (*entities.RefreshTokenClaims, error)
}
