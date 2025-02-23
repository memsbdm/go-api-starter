package ports

import (
	"context"
	"go-starter/internal/domain/entities"
	"time"
)

// TokenService is an interface for interacting with token-related business logic.
type TokenService interface {
	// Generate generates a new token for the given user.
	// Returns the generated token or an error if the generation fails.
	Generate(tokenType entities.TokenType, user *entities.User) (string, error)

	// GenerateTokenWithCache creates a new token for the given user and stores it in cache.
	// Returns the generated token or an error if the operation fails.
	GenerateTokenWithCache(ctx context.Context, tokenType entities.TokenType, user *entities.User) (string, error)

	// ValidateAndParse validates the given token and extracts its claims.
	// Returns a structured representation of the token claims or an error if validation fails.
	ValidateAndParse(tokenType entities.TokenType, token string) (*entities.TokenClaims, error)

	// ValidateAndParseWithCache validates the given token and extracts its claims.
	// Returns a structured representation of the token claims or an error if validation fails or if
	// the token is not stored in cache.
	ValidateAndParseWithCache(ctx context.Context, tokenType entities.TokenType, token string) (*entities.TokenClaims, error)

	// RevokeTokenFromCache deletes the given token from cache.
	// Returns an error if the revocation process fails (e.g., if the token is invalid).
	RevokeTokenFromCache(ctx context.Context, tokenType entities.TokenType, token string) error
}

// TokenProvider is an interface for interacting with token-related data and cryptographic operations.
type TokenProvider interface {
	// Generate generates a new JWT token for the given user.
	// Returns the generated token or an error if the generation fails.
	Generate(tokenType entities.TokenType, user *entities.User, duration time.Duration, signature []byte) (string, error)

	// ValidateAndParse validates the given JWT token and extracts its claims.
	// Returns a structured representation of the token claims or an error if validation fails.
	ValidateAndParse(tokenType entities.TokenType, token string, signature []byte) (*entities.TokenClaims, error)
}
