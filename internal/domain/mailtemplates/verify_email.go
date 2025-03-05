package mailtemplates

import (
	"fmt"
	"time"
)

// VerifyEmail is an email template to validate user's email.
// Returns a string representing the mail body (HTML).
func VerifyEmail(baseURL, token string, expirationTime time.Duration) string {
	return fmt.Sprintf(`Hello, verify your email by visiting <a href="%s/users/me/verify-email/%s">this link</a>!<br><br>This link will expire in %.0f hours.<br>token: %s`, baseURL, token, expirationTime.Hours(), token)
}
