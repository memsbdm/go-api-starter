package mailtemplates

import "fmt"

// VerifyEmail is an email template to validate user's email.
// Returns a string representing the mail body (HTML).
func VerifyEmail(baseURL, token string) string {
	return fmt.Sprintf(`Hello, verify your email by visiting <a href="%s/users/me/email/verify/%s">this link</a>!`, baseURL, token)
}
