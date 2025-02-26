package middleware

import (
	"context"
	"go-starter/internal/adapters/http/helpers"
	"go-starter/internal/adapters/http/responses"
	"go-starter/internal/domain"
	"go-starter/internal/domain/ports"
	"net/http"
)

// AuthMiddleware is a middleware function that validates the authorization token from the incoming HTTP request.
// It sets the user ID in the context of the HTTP request.
func AuthMiddleware(tokenSvc ports.TokenService, errTracker ports.ErrTrackerAdapter) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			accessToken, err := helpers.ExtractTokenFromHeader(r)
			if err != nil {
				responses.HandleError(w, err)
				return
			}

			userID, err := tokenSvc.VerifyAuthToken(r.Context(), accessToken)
			if err != nil {
				responses.HandleError(w, err)
				return
			}

			errTracker.SetUser(userID.String(), r.RemoteAddr)
			ctx := context.WithValue(r.Context(), helpers.AuthorizationPayloadKey, userID.String())
			r = r.WithContext(ctx)

			f(w, r)
		}
	}
}

// GuestMiddleware is a middleware function that allows access to HTTP requests from guests (unauthenticated users).
// It checks for the presence of an authorization token in the request header. If a valid token is found,
// it responds with a forbidden error, preventing authenticated users from accessing guest-only routes.
func GuestMiddleware(tokenSvc ports.TokenService) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			accessToken, err := helpers.ExtractTokenFromHeader(r)
			if err != nil {
				f(w, r)
				return
			}

			_, err = tokenSvc.VerifyAuthToken(r.Context(), accessToken)
			if err == nil {
				responses.HandleError(w, domain.ErrForbidden)
				return
			}

			f(w, r)
		}
	}
}
