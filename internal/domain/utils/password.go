package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes input password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// ComparePassword compares input password with hashed password
func ComparePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
