package http

import (
	"context"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
)

func getAccessTokenClaims(ctx context.Context, key string) (*entities.AccessTokenClaims, error) {
	payload, ok := ctx.Value(key).(*entities.AccessTokenClaims)
	if !ok {
		return nil, domain.ErrTokenClaimsNotFound
	}
	return payload, nil
}
