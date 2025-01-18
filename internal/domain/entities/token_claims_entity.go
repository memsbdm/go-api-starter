package entities

import (
	"github.com/google/uuid"
)

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

// AccessTokenID is a type that represents a unique identifier for an access token, based on UUID.
type AccessTokenID uuid.UUID

// AccessTokenClaims holds the claims associated with an access token.
type AccessTokenClaims struct {
	ID      AccessTokenID
	Subject UserID
}

// UUID converts the AccessTokenID to an uuid.UUID type.
func (id AccessTokenID) UUID() uuid.UUID { return uuid.UUID(id) }

// String returns the string representation of the AccessTokenID.
func (id AccessTokenID) String() string { return uuid.UUID(id).String() }

// ParseAccessTokenID creates a AccessTokenID from a string
func ParseAccessTokenID(s string) (AccessTokenID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return AccessTokenID{}, err
	}
	return AccessTokenID(id), nil
}
