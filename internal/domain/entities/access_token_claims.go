package entities

import "github.com/google/uuid"

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
