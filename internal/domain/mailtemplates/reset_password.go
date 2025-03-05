package mailtemplates

import (
	"fmt"
	"time"
)

// ResetPassword is an email template to reset user's password.
// Returns a string representing the mail body (HTML).
func ResetPassword(baseURL, token string, expirationTime time.Duration) string {
	return fmt.Sprintf(`Hello, reset your password by visiting <a href="%s/auth/password-reset?token=%s">this link</a>!<br><br>This link will expire in %.0f minutes.<br>token: %s`, baseURL, token, expirationTime.Minutes(), token)
}
