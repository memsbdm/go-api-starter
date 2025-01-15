package entities

import "github.com/google/uuid"

// UserID is a type that represents a unique identifier for a user, based on UUID.
type UserID uuid.UUID

// User is an entity that represents a user in the system.
type User struct {
	ID              UserID
	Username        string
	Password        string
	IsEmailVerified bool
}

// UUID converts the UserID to an uuid.UUID type.
func (id UserID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

// String returns the string representation of the UserID.
func (id UserID) String() string {
	return id.UUID().String()
}
