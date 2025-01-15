package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"time"
)

// TokenRepository implements the ports.TokenRepository interface, providing access to
// JWT-related functionalities for token management.
type TokenRepository struct {
	timeGenerator ports.TimeGenerator
}

// NewTokenRepository creates a new instance of TokenRepository.
func NewTokenRepository(timeGenerator ports.TimeGenerator) *TokenRepository {
	return &TokenRepository{
		timeGenerator: timeGenerator,
	}
}

// GenerateAccessToken generates a new JWT access token for the given user.
// Returns the generated access token as a string or an error if the generation fails.
func (tr *TokenRepository) GenerateAccessToken(user *entities.User, duration time.Duration, signature []byte) (string, error) {
	claims := jwt.MapClaims{
		"id":  uuid.New().String(),
		"sub": user.ID.String(),
		"iat": tr.timeGenerator.Now(),
		"exp": tr.timeGenerator.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(signature)
	return signedToken, err
}

// ValidateAndParseAccessToken validates the given JWT access token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails.
func (tr *TokenRepository) ValidateAndParseAccessToken(token string, signature []byte) (*entities.AccessTokenClaims, error) {
	parser := jwt.NewParser(jwt.WithTimeFunc(tr.timeGenerator.Now))

	parsedToken, err := parser.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return signature, nil
	})
	if err != nil {
		return nil, err
	}

	claimsList, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, domain.ErrInvalidToken
	}

	tokenID := claimsList["id"].(string)
	subjectID := claimsList["sub"].(string)

	tokenUUID, err := uuid.Parse(tokenID)
	if err != nil {
		return nil, err
	}

	subjectUUID, err := uuid.Parse(subjectID)
	if err != nil {
		return nil, err
	}

	tokenClaims := &entities.AccessTokenClaims{
		ID:      entities.AccessTokenID(tokenUUID),
		Subject: entities.UserID(subjectUUID),
	}

	return tokenClaims, nil
}

// GenerateRefreshToken creates a new JWT refresh token for the given user ID.
// Returns a unique refresh token ID, the token string, or an error if the operation fails.
func (tr *TokenRepository) GenerateRefreshToken(userID entities.UserID, duration time.Duration, signature []byte) (entities.RefreshTokenID, string, error) {
	id := entities.RefreshTokenID(uuid.New())
	claims := jwt.MapClaims{
		"id":  id.String(),
		"sub": userID.String(),
		"iat": tr.timeGenerator.Now(),
		"exp": tr.timeGenerator.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(signature)
	return id, signedToken, err
}

// ValidateAndParseRefreshToken validates the given JWT refresh token and extracts its claims.
// Returns a structured representation of the token claims or an error if validation fails.
func (tr *TokenRepository) ValidateAndParseRefreshToken(token string, signature []byte) (*entities.RefreshTokenClaims, error) {
	parser := jwt.NewParser(jwt.WithTimeFunc(tr.timeGenerator.Now))

	parsedToken, err := parser.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return signature, nil
	})
	if err != nil {
		return nil, err
	}

	claimsList, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, err
	}

	userID := claimsList["sub"].(string)
	tokenID := claimsList["id"].(string)

	tokenUUID, err := uuid.Parse(tokenID)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	claims := &entities.RefreshTokenClaims{
		ID:      entities.RefreshTokenID(tokenUUID),
		Subject: entities.UserID(userUUID),
	}

	return claims, nil
}