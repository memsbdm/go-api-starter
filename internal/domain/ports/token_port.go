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

	// GenerateOneTimeToken creates a cryptographically secure random token.
	// The token is associated with a user ID and encoded as a base64 string.
	// Returns an error if the token generation fails.
	GenerateOneTimeToken(userID uuid.UUID) (token string, err error)

	// ParseOneTimeToken parses a one-time token and returns the user ID.
	// Returns an error if the token is invalid.
	ParseOneTimeToken(token string) (uuid.UUID, error)
}
