package ports

import (
	"context"
	"go-starter/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

// TokenService is an interface for interacting with token-related business logic.
type TokenService interface {
	// GenerateAccessToken generates a new access token for a user.
	// Returns the signed token string or an error if generation fails.
	GenerateAccessToken(user *entities.User) (string, error)

	// GenerateRefreshToken generates a new refresh token for a user.
	// Returns the signed token string or an error if generation fails.
	GenerateRefreshToken(ctx context.Context, userID entities.UserID) (string, error)

	// VerifyAndParseRefreshToken verifies and parses a refresh token.
	// Returns the parsed token claims or an error if validation fails.
	VerifyAndParseRefreshToken(ctx context.Context, token string) (*entities.RefreshTokenClaims, error)

	// VerifyAndParseAccessToken verifies and parses a JWT access token.
	// Returns the parsed token claims or an error if validation fails.
	VerifyAndParseAccessToken(token string) (*entities.AccessTokenClaims, error)

	// RevokeRefreshToken revokes a refresh token by deleting it from the cache.
	// Returns an error if the token is not found or if the cache deletion fails.
	RevokeRefreshToken(ctx context.Context, token string) error

	// GenerateOneTimeToken generates a new one-time token for a user.
	// Returns the token string or an error if generation fails.
	GenerateOneTimeToken(ctx context.Context, tokenType entities.TokenType, userID entities.UserID) (string, error)

	// VerifyAndConsumeOneTimeToken verifies and consumes a one-time token.
	// Returns the user ID or an error if the token is not found or if the token is invalid.
	VerifyAndConsumeOneTimeToken(ctx context.Context, tokenType entities.TokenType, token string) (entities.UserID, error)
}

// TokenProvider is an interface for interacting with token-related data and cryptographic operations.
type TokenProvider interface {
	// GenerateAccessToken creates a new JWT access token for a user.
	// Returns the signed token string or an error if generation fails.
	GenerateAccessToken(user *entities.User, duration time.Duration, signature []byte) (string, error)

	// GenerateRefreshToken creates a new JWT refresh token for a user.
	// Returns the signed token string or an error if generation fails.
	GenerateRefreshToken(userID uuid.UUID, duration time.Duration, signature []byte) (uuid.UUID, string, error)

	// VerifyAndParseAccessToken verifies and parses a JWT access token.
	// Returns the parsed token claims or an error if validation fails.
	VerifyAndParseAccessToken(accessToken string, signature []byte) (*entities.AccessTokenClaims, error)

	// VerifyAndParseRefreshToken verifies and parses a JWT refresh token.
	// Returns the parsed token claims or an error if validation fails.
	VerifyAndParseRefreshToken(refreshToken string, signature []byte) (*entities.RefreshTokenClaims, error)

	// GenerateOneTimeToken creates a new secure random token associated with a user ID.
	// Returns the token, its hash for storage, and any error that occurred.
	GenerateOneTimeToken(userID uuid.UUID) (token string, hash string, err error)

	// ParseOneTimeToken decodes and validates the structure of a one-time token.
	// Returns the parsed token data or an error if the token is invalid.
	ParseOneTimeToken(token string) (*entities.OneTimeToken, error)

	// HashToken creates a secure hash of the given token for storage and validation.
	// Returns the base64-encoded hash string.
	HashToken(token string) string
}
