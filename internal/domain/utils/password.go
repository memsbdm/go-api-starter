package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the input password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// ComparePassword compares the input password with a hashed password.
// It checks if the provided password matches the hashed password, returning an error if they do not match.
func ComparePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
