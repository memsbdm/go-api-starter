package entities

// TokenType represents the type of token.
type TokenType string

// Token type constants define the available types of tokens in the system.
const (
	AccessToken            TokenType = "access_token"
	EmailVerificationToken TokenType = "email_verification_token"
	PasswordResetToken     TokenType = "password_reset_token"
)

// String converts the TokenType to its string representation.
func (t TokenType) String() string {
	return string(t)
}

// OneTimeToken represents a one-time token.
type OneTimeToken struct {
	UserID UserID
	Token  string
}
