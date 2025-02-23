package helpers

import (
	"context"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
)

const (
	// AuthorizationPayloadKey defines the key used to store and retrieve the authorization payload from the context.
	AuthorizationPayloadKey = "authorization_payload"
)

// ExtractAccessTokenClaims retrieves the AccessTokenClaims from the provided context.
// Returns an error if the claims are not present or cannot be cast to the expected type.
func ExtractAccessTokenClaims(ctx context.Context) (*entities.TokenClaims, error) {
	payload, ok := ctx.Value(AuthorizationPayloadKey).(*entities.TokenClaims)
	if !ok {
		return nil, domain.ErrTokenClaimsNotFound
	}
	return payload, nil
}
