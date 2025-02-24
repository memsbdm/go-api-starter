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

// Provider implements the ports.TokenProvider interface, providing access to
// functionalities for token management.
type Provider struct {
	errTracker    ports.ErrTrackerAdapter
	timeGenerator ports.TimeGenerator
	signingMethod jwt.SigningMethod
}

// NewTokenProvider creates a new instance of JWTProvider.
func NewTokenProvider(timeGenerator ports.TimeGenerator, errTracker ports.ErrTrackerAdapter) *Provider {
	return &Provider{
		timeGenerator: timeGenerator,
		signingMethod: jwt.SigningMethodHS256,
		errTracker:    errTracker,
	}
}

// Generate generates a new JWT token for the given user.
// Returns the generated token or an error if the generation fails.
func (p *Provider) GenerateJWT(tokenType entities.TokenType, user *entities.User, duration time.Duration, signature []byte) (string, error) {
	claims := jwt.MapClaims{
		"id":   uuid.New().String(),
		"sub":  user.ID.String(),
		"iat":  p.timeGenerator.Now().Unix(),
		"exp":  p.timeGenerator.Now().Add(duration).Unix(),
		"type": tokenType.String(),
	}

	token := jwt.NewWithClaims(p.signingMethod, claims)
	signedToken, err := token.SignedString(signature)
	if err != nil {
		p.errTracker.CaptureException(err)
	}
	return signedToken, err
}

// ValidateAndParseJWT validates the given JWT token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails.
func (p *Provider) ValidateAndParseJWT(tokenType entities.TokenType, token string, signature []byte) (*entities.TokenClaims, error) {
	parser := jwt.NewParser(jwt.WithTimeFunc(p.timeGenerator.Now))

	parsedToken, err := parser.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return signature, nil
	})
	if err != nil {
		p.errTracker.CaptureException(err)
		return nil, err
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
	if !ok || entities.TokenType(parsedType) != tokenType {
		return nil, domain.ErrInvalidToken
	}

	claims := &entities.TokenClaims{
		ID:      tokenUUID,
		Subject: userID,
		Type:    tokenType,
	}

	return claims, nil
}

// GenerateSecureToken creates a new secure random token associated with a user ID.
// Returns the token, its hash for storage, and any error that occurred.
func (p *Provider) GenerateSecureToken(userID uuid.UUID) (token string, hash string, err error) {
	// Générer bytes aléatoires
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		p.errTracker.CaptureException(err)
		return "", "", err
	}

	// Créer le token avec userID
	randomToken := &entities.SecureToken{
		UserID: userID,
		Token:  base64.RawURLEncoding.EncodeToString(b),
	}

	// Encoder en JSON
	tokenJSON, err := json.Marshal(randomToken)
	if err != nil {
		p.errTracker.CaptureException(err)
		return "", "", err
	}

	// Encoder en base64
	token = base64.RawURLEncoding.EncodeToString(tokenJSON)
	// Générer le hash
	hash = p.HashSecureToken(token)

	return token, hash, nil
}

// HashSecureToken creates a secure hash of the given token for storage and validation.
// Returns the base64-encoded hash string.
func (p *Provider) HashSecureToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// ParseSecureToken decodes and validates the structure of a secure token.
// Returns the parsed token data or an error if the token is invalid.
func (p *Provider) ParseSecureToken(token string) (*entities.SecureToken, error) {
	// Decoder le base64
	tokenJSON, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		p.errTracker.CaptureException(err)
		return nil, fmt.Errorf("invalid token encoding: %w", err)
	}

	// Parser le JSON
	var randomToken entities.SecureToken
	if err := json.Unmarshal(tokenJSON, &randomToken); err != nil {
		p.errTracker.CaptureException(err)
		return nil, fmt.Errorf("invalid token format: %w", err)
	}

	return &randomToken, nil
}
