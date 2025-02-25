package entities

import (
	"github.com/google/uuid"
)

// TokenType represents the type of token.
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

// AccessTokenClaims represents the claims of an access token.
type AccessTokenClaims struct {
	ID      uuid.UUID
	Subject UserID
	Type    TokenType
}

// RefreshTokenClaims represents the claims of a refresh token.
type RefreshTokenClaims struct {
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

// OneTimeToken represents a one-time token.
type OneTimeToken struct {
	UserID UserID
	Token  string
}
