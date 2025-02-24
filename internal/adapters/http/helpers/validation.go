package helpers

import (
	"go-starter/internal/domain"
	"regexp"
)

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.+_-]+@[a-zA-Z0-9.+_-]+$")
	if len(email) > domain.EmailMaxLength || !emailRegex.MatchString(email) {
		return false
	}
	return true
}
