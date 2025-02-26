package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Provider implements the ports.TokenProvider interface, providing access to functionalities for token management.
// It manages one-time and authentication tokens.
type Provider struct {
	errTracker    ports.ErrTrackerAdapter
	timeGenerator ports.TimeGenerator
	signingMethod jwt.SigningMethod
}

// NewTokenProvider creates a new instance of Provider.
func NewTokenProvider(timeGenerator ports.TimeGenerator, errTracker ports.ErrTrackerAdapter) *Provider {
	return &Provider{
		timeGenerator: timeGenerator,
		signingMethod: jwt.SigningMethodHS256,
		errTracker:    errTracker,
	}
}

// GenerateRandomToken creates a cryptographically secure random token.
// The token is encoded as a base64 string.
// Returns:
// - token: the secure token to be sent to the user
// - error: any error that occurred during generation
func (p *Provider) GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateOneTimeToken creates a cryptographically secure random token.
// The token is associated with a user ID and encoded as a base64 string.
// Returns:
// - token: the secure token to be sent to the user
// - hash: the hashed version of the token for storage
// - error: any error that occurred during generation
func (p *Provider) GenerateOneTimeToken(userID uuid.UUID) (token string, hash string, err error) {
	// Generate 32 bytes of random data (256 bits)
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		p.errTracker.CaptureException(fmt.Errorf("failed to generate secure token: %w", err))
		return "", "", err
	}

	oneTimeToken := &entities.OneTimeToken{
		UserID: entities.UserID(userID),
		Token:  base64.RawURLEncoding.EncodeToString(randomBytes),
	}

	tokenJSON, err := json.Marshal(oneTimeToken)
	if err != nil {
		p.errTracker.CaptureException(fmt.Errorf("failed to marshal secure token: %w", err))
		return "", "", err
	}

	token = base64.RawURLEncoding.EncodeToString(tokenJSON)
	hash = p.HashToken(token)

	return token, hash, nil
}

// ParseOneTimeToken decodes and validates the structure of a one-time token.
// Returns the parsed token data or an error if the token is invalid.
func (p *Provider) ParseOneTimeToken(token string) (*entities.OneTimeToken, error) {
	tokenJSON, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		err = fmt.Errorf("invalid token encoding: %w", err)
		p.errTracker.CaptureException(err)
		return nil, err
	}

	var oneTimeToken entities.OneTimeToken
	if err := json.Unmarshal(tokenJSON, &oneTimeToken); err != nil {
		err = fmt.Errorf("invalid token format: %w", err)
		p.errTracker.CaptureException(err)
		return nil, err
	}

	return &oneTimeToken, nil
}

// HashToken creates a secure hash of the given token for storage and validation.
// Returns the base64-encoded hash string.
func (p *Provider) HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(h[:])
}
