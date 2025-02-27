package helpers

import (
	"go-starter/internal/domain"
	"regexp"
)

// IsValidEmail checks if an email is valid.
// Returns true if the email is valid, false otherwise.
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.+_-]+@[a-zA-Z0-9.+_-]+$")
	if len(email) > domain.EmailMaxLength || !emailRegex.MatchString(email) {
		return false
	}
	return true
}
