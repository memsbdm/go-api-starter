package mailtemplates

import "fmt"

// VerifyEmail is an email template to validate user's email.
// Returns a string representing the mail body (HTML).
func VerifyEmail(baseURL, token string) string {
	return fmt.Sprintf(`Hello, verify your email by visiting <a href="%s/users/me/verify-email/%s">this link</a>!<br><br>This link will expire in 24 hours.<br>token: %s`, baseURL, token, token)
}
