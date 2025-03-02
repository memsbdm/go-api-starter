package mailtemplates

import "fmt"

// ResetPassword is an email template to reset user's password.
// Returns a string representing the mail body (HTML).
func ResetPassword(baseURL, token string) string {
	return fmt.Sprintf(`Hello, reset your password by visiting <a href="%s/users/me/password/reset/%s">this link</a>!<br><br>This link will expire in 15 minutes.<br>token: %s`, baseURL, token, token)
}
