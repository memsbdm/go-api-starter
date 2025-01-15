package entities

import "github.com/google/uuid"

// RefreshTokenID is a type that represents a unique identifier for a refresh token, based on UUID.
type RefreshTokenID uuid.UUID

// RefreshTokenClaims holds the claims associated with a refresh token.
type RefreshTokenClaims struct {
	ID      RefreshTokenID
	Subject UserID
}

// UUID converts the RefreshTokenID to an uuid.UUID type.
func (id RefreshTokenID) UUID() uuid.UUID { return uuid.UUID(id) }

// String returns the string representation of the RefreshTokenID.
func (id RefreshTokenID) String() string { return uuid.UUID(id).String() }
