package entities

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

// TokenPayload represents authentication claims
type TokenPayload struct {
	ID        uuid.UUID
	UserID    UserID
	ExpiresAt int64
	IssuedAt  int64
}

// GetIssuer returns the issuer of the token, which is typically the service that generated it
func (tp *TokenPayload) GetIssuer() (string, error) {
	return "", nil
}

// GetSubject returns the subject of the token, typically representing the entity the token is issued to
func (tp *TokenPayload) GetSubject() (string, error) {
	return "", nil
}

// GetAudience returns the audience for which the token is intended
func (tp *TokenPayload) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

// GetExpirationTime returns the expiration time of the token as a NumericDate
func (tp *TokenPayload) GetExpirationTime() (*jwt.NumericDate, error) {
	if tp.ExpiresAt == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(tp.ExpiresAt, 0)), nil
}

// GetIssuedAt returns the issuance time of the token as a NumericDate
func (tp *TokenPayload) GetIssuedAt() (*jwt.NumericDate, error) {
	if tp.IssuedAt == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(tp.IssuedAt, 0)), nil
}

// GetNotBefore returns the not-before time for the token as a NumericDate
func (tp *TokenPayload) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

// Valid checks if the token is valid based on its expiration time
func (tp *TokenPayload) Valid() error {
	if tp.ExpiresAt != 0 && time.Now().Unix() > tp.ExpiresAt {
		return fmt.Errorf("token has expired")
	}
	return nil
}
