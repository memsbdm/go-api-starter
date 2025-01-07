package entities

import "github.com/google/uuid"

type UserID uuid.UUID

// User is an entity that represents a user
type User struct {
	ID       UserID `json:"id"`
	Username string `json:"username"`
}

func (id UserID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id UserID) String() string {
	return id.UUID().String()
}
