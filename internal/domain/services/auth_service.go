package services

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go-starter/config"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/utils"
	"time"
)

// AuthService implements ports.AuthService interface
type AuthService struct {
	cfg     config.Security
	userSvc ports.UserService
}

// NewAuthService creates a auth service instance
func NewAuthService(cfg *config.Security, userSvc ports.UserService) *AuthService {
	return &AuthService{
		cfg:     *cfg,
		userSvc: userSvc,
	}
}

// Login authenticates a user. Returns a token string pointer (could be JWT, session ID, etc.) or an error if login fails
func (as *AuthService) Login(ctx context.Context, username, password string) (*string, error) {
	user, err := as.userSvc.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	err = utils.ComparePassword(password, user.Password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	jtwStr, err := GenerateJWT(as.cfg.JWTSecret, user.ID)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return &jtwStr, nil
}

// GenerateJWT generates a JWT
func GenerateJWT(jwtSecretKey []byte, userID entities.UserID) (string, error) {
	claims := entities.TokenPayload{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

// ValidateToken validates an auth token (JWT, session ID, etc.) and returns associated claims
func (as *AuthService) ValidateToken(tokenStr string) (*entities.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &entities.TokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		return as.cfg.JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}

	tokenPayload, ok := token.Claims.(*entities.TokenPayload)
	if !ok || !token.Valid {
		return nil, domain.ErrUnauthorized
	}

	return tokenPayload, nil
}
