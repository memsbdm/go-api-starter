package entities

import "github.com/google/uuid"

type AccessTokenClaims struct {
	ID      uuid.UUID
	Subject UserID
}
