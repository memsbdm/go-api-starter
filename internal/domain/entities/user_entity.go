package entities

import (
	"go-starter/internal/domain"
	"time"

	"github.com/google/uuid"
)

// UserID is a type that represents a unique identifier for a user, based on UUID.
type UserID uuid.UUID

// User is an entity that represents a user in the system.
type User struct {
	ID              UserID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Name            string
	Username        string
	Password        string
	Email           string
	IsEmailVerified bool
	RoleID          RoleID
	AvatarURL       *string
}

// NilUserID is the nil UserID.
var NilUserID = UserID(uuid.Nil)

// UUID converts the UserID to an uuid.UUID type.
func (id UserID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

// String returns the string representation of the UserID.
func (id UserID) String() string {
	return id.UUID().String()
}

// ParseUserID creates a UserID from a string
func ParseUserID(s string) (UserID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UserID{}, domain.ErrInvalidUserId
	}
	return UserID(id), nil
}

// RoleID is a type that represents a role ID for a user, based on int.
type RoleID int

const (
	RoleAdmin RoleID = iota
	RoleUser
)

// Int returns the integer representation of the RoleID.
func (r RoleID) Int() int {
	return int(r)
}

// UpdateUserParams holds the parameters required for updating a user's information.
type UpdateUserParams struct {
	Password             *string
	PasswordConfirmation *string
}
