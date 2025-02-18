package mailtemplates

import "fmt"

// Hello is an example of email template.
// Returns a string representing the mail body (HTML).
func Hello(name string) string {
	return fmt.Sprintf("Hello, %s!<br>Nice to meet you!", name)
}
