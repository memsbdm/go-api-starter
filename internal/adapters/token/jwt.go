package token

import (
	"fmt"
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

// GenerateAccessToken generates a new JWT access token for the given user.
// Returns the generated access token or an error if the generation fails.
func (jp *JWTProvider) GenerateAccessToken(user *entities.User, duration time.Duration, signature []byte) (entities.AccessToken, error) {
	claims := jwt.MapClaims{
		"id":  uuid.New().String(),
		"sub": user.ID.String(),
		"iat": jp.timeGenerator.Now().Unix(),
		"exp": jp.timeGenerator.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jp.signingMethod, claims)
	signedToken, err := token.SignedString(signature)
	if err != nil {
		jp.errTracker.CaptureException(err)
	}
	return entities.AccessToken(signedToken), err
}

// ValidateAndParseAccessToken validates the given JWT access token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails.
func (jp *JWTProvider) ValidateAndParseAccessToken(token string, signature []byte) (*entities.AccessTokenClaims, error) {
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

	sub, ok := claimsList["sub"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	tokenUUID, err := entities.ParseAccessTokenID(id)
	if err != nil {
		err = fmt.Errorf("could not parse access token %s: %w", id, err)
		jp.errTracker.CaptureException(err)
		return nil, err
	}

	userID, err := entities.ParseUserID(sub)
	if err != nil {
		err = fmt.Errorf("could not parse user id %s: %w", sub, err)
		jp.errTracker.CaptureException(err)
		return nil, err
	}

	tokenClaims := &entities.AccessTokenClaims{
		ID:      tokenUUID,
		Subject: userID,
	}

	return tokenClaims, nil
}

// GenerateRefreshToken creates a new JWT refresh token for the given user ID.
// Returns a unique refresh token ID, the refresh token, or an error if the operation fails.
func (jp *JWTProvider) GenerateRefreshToken(userID entities.UserID, duration time.Duration, signature []byte) (entities.RefreshTokenID, entities.RefreshToken, error) {
	id := entities.RefreshTokenID(uuid.New())
	claims := jwt.MapClaims{
		"id":  id.String(),
		"sub": userID.String(),
		"iat": jp.timeGenerator.Now().Unix(),
		"exp": jp.timeGenerator.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jp.signingMethod, claims)
	signedToken, err := token.SignedString(signature)
	if err != nil {
		jp.errTracker.CaptureException(err)
	}
	return id, entities.RefreshToken(signedToken), err
}

// ValidateAndParseRefreshToken validates the given JWT refresh token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails.
func (jp *JWTProvider) ValidateAndParseRefreshToken(token string, signature []byte) (*entities.RefreshTokenClaims, error) {
	parser := jwt.NewParser(jwt.WithTimeFunc(jp.timeGenerator.Now))

	parsedToken, err := parser.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return signature, nil
	})
	if err != nil {
		err = fmt.Errorf("could not parse refresh token %s: %w", token, err)
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

	sub, ok := claimsList["sub"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	tokenUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("could not parse refresh token %s: %w", id, err)
		jp.errTracker.CaptureException(err)
		return nil, err
	}

	userID, err := entities.ParseUserID(sub)
	if err != nil {
		err = fmt.Errorf("could not parse user id %s: %w", sub, err)
		jp.errTracker.CaptureException(err)
		return nil, err
	}

	return &entities.RefreshTokenClaims{
		ID:      entities.RefreshTokenID(tokenUUID),
		Subject: userID,
	}, nil
}
