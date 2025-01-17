package responses

import (
	"go-starter/internal/domain/entities"
	"go-starter/pkg/env"
)

// LoginResponse represents the structure of a response body for successful authentication.
type LoginResponse struct {
	Tokens AuthTokensResponse `json:"tokens"`
	User   UserResponse       `json:"user"`
}

// AuthTokensResponse represents the structure of a response body with an access token, a refresh token
// and their expiration time in ms.
type AuthTokensResponse struct {
	AccessToken             string `json:"access_token" example:"eyJhbGciOi..."`
	RefreshToken            string `json:"refresh_token" example:"eyJhbGciOi..."`
	AccessTokenExpiredInMs  int64  `json:"access_token_expired_in_ms" example:"50000"`
	RefreshTokenExpiredInMs int64  `json:"refresh_token_expired_in_ms" example:"890000"`
}

// NewAuthTokensResponse is a helper function that creates a AuthTokensResponse from an auth tokens entity.
func NewAuthTokensResponse(tokens *entities.AuthTokens) AuthTokensResponse {
	accessTokenDuration := env.GetDuration("ACCESS_TOKEN_DURATION")
	refreshTokenDuration := env.GetDuration("REFRESH_TOKEN_DURATION")
	return AuthTokensResponse{
		AccessToken:             string(tokens.AccessToken),
		RefreshToken:            string(tokens.RefreshToken),
		AccessTokenExpiredInMs:  accessTokenDuration.Milliseconds(),
		RefreshTokenExpiredInMs: refreshTokenDuration.Milliseconds(),
	}
}

// NewLoginResponse is a helper function that creates a LoginResponse with the provided auth tokens and user.
func NewLoginResponse(tokens *entities.AuthTokens, user *entities.User) LoginResponse {
	return LoginResponse{
		Tokens: NewAuthTokensResponse(tokens),
		User:   NewUserResponse(user),
	}
}

// RefreshTokenResponse represents the structure of a response body for successful token refresh.
type RefreshTokenResponse struct {
	Tokens AuthTokensResponse `json:"tokens"`
}

// NewRefreshTokenResponse is a helper function that creates a RefreshTokenResponse with the provided auth tokens.
func NewRefreshTokenResponse(tokens *entities.AuthTokens) RefreshTokenResponse {
	return RefreshTokenResponse{
		Tokens: NewAuthTokensResponse(tokens),
	}
}
