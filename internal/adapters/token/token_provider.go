package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"go-starter/internal/domain/ports"
	"strings"

	"github.com/google/uuid"
)

// Provider implements the ports.TokenProvider interface, providing access to functionalities for token management.
// It manages one-time and authentication tokens.
type Provider struct {
	errTracker    ports.ErrTrackerAdapter
	timeGenerator ports.TimeGenerator
}

// NewTokenProvider creates a new instance of Provider.
func NewTokenProvider(timeGenerator ports.TimeGenerator, errTracker ports.ErrTrackerAdapter) *Provider {
	return &Provider{
		timeGenerator: timeGenerator,
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
		p.errTracker.CaptureException(err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateOneTimeToken creates a cryptographically secure random token.
// The token is associated with a user ID and encoded as a base64 string.
// Returns an error if the token generation fails.
func (p *Provider) GenerateOneTimeToken(userID uuid.UUID) (token string, err error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		p.errTracker.CaptureException(err)
		return "", err
	}
	randomPart := base64.URLEncoding.EncodeToString(randomBytes)
	composite := userID.String() + "." + randomPart

	return base64.URLEncoding.EncodeToString([]byte(composite)), nil
}

// ParseOneTimeToken parses a one-time token and returns the user ID.
// Returns an error if the token is invalid.
func (p *Provider) ParseOneTimeToken(token string) (uuid.UUID, error) {
	decodedBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return uuid.Nil, err
	}

	decoded := string(decodedBytes)
	parts := strings.Split(decoded, ".")
	if len(parts) != 2 {
		return uuid.Nil, errors.New("invalid one-time token format")
	}

	userID, err := uuid.Parse(parts[0])
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
