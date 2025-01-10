package entities

import "github.com/golang-jwt/jwt/v5"

// TokenPayload represents authentication claims
type TokenPayload struct {
	UserID UserID
	jwt.RegisteredClaims
}
