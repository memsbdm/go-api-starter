package http

import (
	"context"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
)

func getAuthPayload(ctx context.Context, key string) (*entities.TokenPayload, error) {
	payload, ok := ctx.Value(key).(*entities.TokenPayload)
	if !ok {
		return nil, domain.ErrTokenPayloadNotFound
	}
	return payload, nil
}
