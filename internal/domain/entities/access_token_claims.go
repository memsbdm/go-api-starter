package entities

import "github.com/google/uuid"

type AccessTokenID uuid.UUID

type AccessTokenClaims struct {
	ID      AccessTokenID
	Subject UserID
}

func (id AccessTokenID) UUID() uuid.UUID { return uuid.UUID(id) }

func (id AccessTokenID) String() string { return uuid.UUID(id).String() }
