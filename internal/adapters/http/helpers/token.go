package helpers

import (
	"context"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"net/http"
	"strings"
)

const (
	// AuthorizationHeaderKey defines the key used to retrieve the authorization header from the HTTP request.
	AuthorizationHeaderKey = "Authorization"
	// AuthorizationType specifies the accepted type of authorization.
	AuthorizationType = "bearer"
	// AuthorizationPayloadKey defines the key used to store and retrieve the authorization payload from the context.
	AuthorizationPayloadKey = "authorization_payload"
)

// ExtractTokenFromHeader extracts the token from the authorization header of the HTTP request.
// Returns the token or an error if the token is not found or if the token is invalid.
func ExtractTokenFromHeader(r *http.Request) (string, error) {
	authorizationHeader := r.Header.Get(AuthorizationHeaderKey)
	if len(authorizationHeader) == 0 {
		return "", domain.ErrUnauthorized
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) != 2 {
		return "", domain.ErrUnauthorized
	}

	if !strings.EqualFold(fields[0], AuthorizationType) {
		return "", domain.ErrUnauthorized
	}

	return fields[1], nil
}

// GetUserIDFromContext retrieves the user ID from the context of the HTTP request.
// Returns the user ID or an error if the user ID is not found or if the user ID is invalid.
func GetUserIDFromContext(ctx context.Context) (entities.UserID, error) {
	id, ok := ctx.Value(AuthorizationPayloadKey).(string)
	if !ok {
		return entities.NilUserID, domain.ErrInternal
	}

	userID, err := entities.ParseUserID(id)
	if err != nil {
		return userID, domain.ErrInternal
	}

	return userID, nil
}
