package entities

import "github.com/google/uuid"

type RefreshTokenClaims struct {
	ID      uuid.UUID
	Subject UserID
}
