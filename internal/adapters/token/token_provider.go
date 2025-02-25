package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Provider implements the ports.TokenProvider interface, providing access to functionalities for token management.
// It handles both JWT tokens for authentication and secure tokens for email verification, password reset etc.
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

// GenerateAccessToken creates a new JWT access token for a user.
// Returns the signed token string or an error if generation fails.
func (p *Provider) GenerateAccessToken(user *entities.User, duration time.Duration, signature []byte) (string, error) {
	tokenID := uuid.New()
	claims := jwt.MapClaims{
		"id":   tokenID.String(),
		"sub":  user.ID.String(),
		"iat":  p.timeGenerator.Now().Unix(),
		"exp":  p.timeGenerator.Now().Add(duration).Unix(),
		"type": entities.AccessToken,
	}
	token := jwt.NewWithClaims(p.signingMethod, claims)
	signedToken, err := token.SignedString(signature)
	if err != nil {
		p.errTracker.CaptureException(fmt.Errorf("failed to sign access token: %w", err))
		return "", err
	}
	return signedToken, nil
}

// GenerateRefreshToken creates a new JWT refresh token for a user.
// Returns the signed token string or an error if generation fails.
func (p *Provider) GenerateRefreshToken(userID uuid.UUID, duration time.Duration, signature []byte) (uuid.UUID, string, error) {
	tokenID := uuid.New()
	claims := jwt.MapClaims{
		"id":   tokenID.String(),
		"sub":  userID.String(),
		"iat":  p.timeGenerator.Now().Unix(),
		"exp":  p.timeGenerator.Now().Add(duration).Unix(),
		"type": entities.RefreshToken,
	}
	token := jwt.NewWithClaims(p.signingMethod, claims)
	signedToken, err := token.SignedString(signature)
	if err != nil {
		p.errTracker.CaptureException(fmt.Errorf("failed to sign access token: %w", err))
		return uuid.Nil, "", err
	}
	return tokenID, signedToken, nil
}

// VerifyAndParseAccessToken verifies and parses a JWT access token.
// Returns the parsed token claims or an error if validation fails.
func (p *Provider) VerifyAndParseAccessToken(accessToken string, signature []byte) (*entities.AccessTokenClaims, error) {
	parser := jwt.NewParser(
		jwt.WithTimeFunc(p.timeGenerator.Now),
		jwt.WithValidMethods([]string{p.signingMethod.Alg()}),
	)

	parsedToken, err := parser.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return signature, nil
	})
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claimsList, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, domain.ErrInvalidToken
	}

	id, ok := claimsList["id"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	tokenUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	sub, ok := claimsList["sub"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	userID, err := entities.ParseUserID(sub)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	parsedType, ok := claimsList["type"].(string)
	if !ok || entities.TokenType(parsedType) != entities.AccessToken {
		return nil, domain.ErrInvalidToken
	}

	return &entities.AccessTokenClaims{
		ID:      tokenUUID,
		Subject: userID,
		Type:    entities.AccessToken,
	}, nil
}

// VerifyAndParseRefreshToken verifies and parses a JWT refresh token.
// Returns the parsed token claims or an error if validation fails.
func (p *Provider) VerifyAndParseRefreshToken(refreshToken string, signature []byte) (*entities.RefreshTokenClaims, error) {
	parser := jwt.NewParser(
		jwt.WithTimeFunc(p.timeGenerator.Now),
		jwt.WithValidMethods([]string{p.signingMethod.Alg()}),
	)

	parsedToken, err := parser.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return signature, nil
	})
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claimsList, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, domain.ErrInvalidToken
	}

	id, ok := claimsList["id"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	tokenUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	sub, ok := claimsList["sub"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	userID, err := entities.ParseUserID(sub)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	parsedType, ok := claimsList["type"].(string)
	if !ok || entities.TokenType(parsedType) != entities.RefreshToken {
		return nil, domain.ErrInvalidToken
	}

	return &entities.RefreshTokenClaims{
		ID:      tokenUUID,
		Subject: userID,
		Type:    entities.RefreshToken,
	}, nil
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
	hash = p.HashOneTimeToken(token)

	return token, hash, nil
}

// HashOneTimeToken creates a secure hash of the given token for storage and validation.
// Returns the base64-encoded hash string.
func (p *Provider) HashOneTimeToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(h[:])
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
