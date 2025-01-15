package entities

import "github.com/google/uuid"

type RefreshTokenID uuid.UUID

type RefreshTokenClaims struct {
	ID      RefreshTokenID
	Subject UserID
}

func (id RefreshTokenID) UUID() uuid.UUID { return uuid.UUID(id) }

func (id RefreshTokenID) String() string { return uuid.UUID(id).String() }
