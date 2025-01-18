package entities

// AccessToken is a type that represents an access token, based on string.
type AccessToken string

// RefreshToken is a type that represents a refresh token, based on string.
type RefreshToken string

// AuthTokens represents a pair of authentication tokens containing both an access token and a refresh token.
type AuthTokens struct {
	AccessToken  AccessToken
	RefreshToken RefreshToken
}
