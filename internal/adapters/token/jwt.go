package token

import (
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTProvider implements the ports.TokenProvider interface, providing access to
// JWT-related functionalities for token management.
type JWTProvider struct {
	errTracker    ports.ErrTrackerAdapter
	timeGenerator ports.TimeGenerator
	signingMethod jwt.SigningMethod
}

// NewJWTProvider creates a new instance of JWTProvider.
func NewJWTProvider(timeGenerator ports.TimeGenerator, errTracker ports.ErrTrackerAdapter) *JWTProvider {
	return &JWTProvider{
		timeGenerator: timeGenerator,
		signingMethod: jwt.SigningMethodHS256,
		errTracker:    errTracker,
	}
}

// Generate generates a new JWT token for the given user.
// Returns the generated token or an error if the generation fails.
func (jp *JWTProvider) Generate(tokenType entities.TokenType, user *entities.User, duration time.Duration, signature []byte) (string, error) {
	claims := jwt.MapClaims{
		"id":   uuid.New().String(),
		"sub":  user.ID.String(),
		"iat":  jp.timeGenerator.Now().Unix(),
		"exp":  jp.timeGenerator.Now().Add(duration).Unix(),
		"type": tokenType.String(),
	}

	token := jwt.NewWithClaims(jp.signingMethod, claims)
	signedToken, err := token.SignedString(signature)
	if err != nil {
		jp.errTracker.CaptureException(err)
	}
	return signedToken, err
}

// ValidateAndParse validates the given JWT token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails.
func (jp *JWTProvider) ValidateAndParse(tokenType entities.TokenType, token string, signature []byte) (*entities.TokenClaims, error) {
	parser := jwt.NewParser(jwt.WithTimeFunc(jp.timeGenerator.Now))

	parsedToken, err := parser.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return signature, nil
	})
	if err != nil {
		jp.errTracker.CaptureException(err)
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
