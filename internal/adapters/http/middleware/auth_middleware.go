package middleware

import (
	"context"
	"go-starter/internal/adapters/http/helpers"
	"go-starter/internal/adapters/http/responses"
	"go-starter/internal/domain"
	"go-starter/internal/domain/ports"
	"net/http"
	"strings"
)

const (
	// authorizationHeaderKey defines the key used to retrieve the authorization header from the HTTP request.
	authorizationHeaderKey = "Authorization"
	// authorizationType specifies the accepted type of authorization.
	authorizationType = "bearer"
)

// AuthMiddleware is a middleware function that validates the authorization token from the incoming HTTP request.
func AuthMiddleware(tokenService *ports.TokenService) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get(authorizationHeaderKey)
			if len(authorizationHeader) == 0 {
				responses.HandleError(w, domain.ErrUnauthorized)
				return
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) != 2 {
				responses.HandleError(w, domain.ErrUnauthorized)
				return
			}

			if strings.ToLower(fields[0]) != strings.ToLower(authorizationType) {
				responses.HandleError(w, domain.ErrUnauthorized)
				return
			}

			accessToken := fields[1]
			tokenPayload, err := (*tokenService).ValidateAndParseAccessToken(accessToken)
			if err != nil {
				responses.HandleError(w, domain.ErrUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), helpers.AuthorizationPayloadKey, tokenPayload)
			r = r.WithContext(ctx)

			f(w, r)
		}
	}
}

// GuestMiddleware is a middleware function that allows access to HTTP requests from guests (unauthenticated users).
// It checks for the presence of an authorization token in the request header. If a valid token is found,
// it responds with a forbidden error, preventing authenticated users from accessing guest-only routes.
func GuestMiddleware(tokenService *ports.TokenService) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get(authorizationHeaderKey)
			if len(authorizationHeader) == 0 {
				f(w, r)
				return
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) != 2 || strings.ToLower(fields[0]) != strings.ToLower(authorizationType) {
				f(w, r)
				return
			}

			accessToken := fields[1]
			_, err := (*tokenService).ValidateAndParseAccessToken(accessToken)
			if err == nil {
				responses.HandleError(w, domain.ErrForbidden)
				return
			}

			f(w, r)
		}
	}
}
