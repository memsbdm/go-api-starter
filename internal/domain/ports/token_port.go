package ports

import (
	"context"
	"go-starter/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

// TokenService is an interface for interacting with token-related business logic.
type TokenService interface {
	// GenerateJWT generates a new token for the given user.
	// Returns the generated token or an error if the generation fails.
	GenerateJWT(tokenType entities.TokenType, user *entities.User) (string, error)

	// CreateAndCacheJWT creates a new token for the given user and stores it in cache.
	// Returns the generated token or an error if the operation fails.
	CreateAndCacheJWT(ctx context.Context, tokenType entities.TokenType, user *entities.User) (string, error)

	// CreateAndCacheSecureToken creates a new secure token for the given user and stores it in cache.
	// Returns the generated token or an error if the operation fails.
	CreateAndCacheSecureToken(ctx context.Context, tokenType entities.TokenType, user *entities.User) (string, error)

	// VerifyAndInvalidateSecureToken validates the secure token and removes it from cache.
	// Returns the user ID associated with the token or an error if validation fails.
	VerifyAndInvalidateSecureToken(ctx context.Context, tokenType entities.TokenType, token string) (uuid.UUID, error)

	// ValidateJWT validates the given token and extracts its claims.
	// Returns a structured representation of the token claims or an error if validation fails.
	ValidateJWT(tokenType entities.TokenType, token string) (*entities.TokenClaims, error)

	// VerifyCachedJWT validates the given token and verifies its presence in cache.
	// Returns a structured representation of the token claims or an error if validation fails.
	VerifyCachedJWT(ctx context.Context, tokenType entities.TokenType, token string) (*entities.TokenClaims, error)

	// RevokeJWT deletes the given token from cache.
	// Returns an error if the revocation process fails.
	RevokeJWT(ctx context.Context, tokenType entities.TokenType, token string) error
}

// TokenProvider is an interface for interacting with token-related data and cryptographic operations.
type TokenProvider interface {
	// GenerateJWT generates a new JWT token for the given user with specified duration and signature.
	// Returns the generated token or an error if the generation fails.
	GenerateJWT(tokenType entities.TokenType, user *entities.User, duration time.Duration, signature []byte) (string, error)

	// ValidateAndParseJWT validates the given JWT token and extracts its claims.
	// Returns a structured representation of the token claims or an error if validation fails.
	ValidateAndParseJWT(tokenType entities.TokenType, token string, signature []byte) (*entities.TokenClaims, error)

	// GenerateSecureToken creates a new secure random token associated with a user ID.
	// Returns the token, its hash for storage, and any error that occurred.
	GenerateSecureToken(userID uuid.UUID) (token string, hash string, err error)

	// HashSecureToken creates a secure hash of the given token for storage and validation.
	// Returns the base64-encoded hash string.
	HashSecureToken(token string) string

	// ParseSecureToken decodes and validates the structure of a secure token.
	// Returns the parsed token data or an error if the token is invalid.
	ParseSecureToken(token string) (*entities.SecureToken, error)
}
