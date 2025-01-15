package http

import (
	"context"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
)

// ExtractAccessTokenClaims retrieves the AccessTokenClaims from the provided context.
// Returns an error if the claims are not present or cannot be cast to the expected type.
func extractAccessTokenClaims(ctx context.Context) (*entities.AccessTokenClaims, error) {
	payload, ok := ctx.Value(authorizationPayloadKey).(*entities.AccessTokenClaims)
	if !ok {
		return nil, domain.ErrTokenClaimsNotFound
	}
	return payload, nil
}
