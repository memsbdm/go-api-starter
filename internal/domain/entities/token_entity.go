package entities

import (
	"github.com/google/uuid"
)

// TokenType represents the type of the token.
type TokenType string

// Token type constants define the available types of tokens in the system.
const (
	RefreshToken           TokenType = "refresh_token"
	AccessToken            TokenType = "access_token"
	EmailVerificationToken TokenType = "email_verification_token"
)

// String converts the TokenType to its string representation.
func (t TokenType) String() string {
	return string(t)
}

// TokenClaims holds the claims associated with a token.
type TokenClaims struct {
	ID      uuid.UUID
	Subject UserID
	Type    TokenType
}

// AuthTokens contains both access and refresh tokens.
// This structure is typically used when returning authentication credentials to clients.
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

// SecureToken contains a user ID and a token.
// This structure is typically used when returning a secure token to clients and parsing it back.
type SecureToken struct {
	UserID uuid.UUID
	Token  string
}
