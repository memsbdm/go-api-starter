package ports

import (
	"context"
	"go-starter/internal/domain/entities"

	"github.com/google/uuid"
)

// TokenService is an interface for interacting with token-related business logic.
type TokenService interface {
	// GenerateAuthToken generates an access token for a user.
	// Returns the access token or an error if generation fails.
	GenerateAuthToken(ctx context.Context, userID entities.UserID) (string, error)

	// VerifyAuthToken verifies an access token.
	// Returns the user ID or an error if the token is not found or if the token is invalid.
	VerifyAuthToken(ctx context.Context, token string) (entities.UserID, error)

	// RevokeAuthToken revokes an access token.
	// Returns an error if the revocation fails.
	RevokeAuthToken(ctx context.Context, token string) error

	// GenerateOneTimeToken generates a new one-time token for a user.
	// Returns the token string or an error if generation fails.
	GenerateOneTimeToken(ctx context.Context, tokenType entities.TokenType, userID entities.UserID) (string, error)

	// VerifyAndConsumeOneTimeToken verifies and consumes a one-time token.
	// Returns the user ID or an error if the token is not found or if the token is invalid.
	VerifyAndConsumeOneTimeToken(ctx context.Context, tokenType entities.TokenType, token string) (entities.UserID, error)
}

// TokenProvider is an interface for interacting with token-related data and cryptographic operations.
type TokenProvider interface {
	// GenerateRandomToken creates a new secure random token.
	// Returns the token or an error if the operation fails.
	GenerateRandomToken() (string, error)

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
