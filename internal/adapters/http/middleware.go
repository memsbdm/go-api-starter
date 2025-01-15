package http

import (
	"context"
	"go-starter/internal/domain"
	"go-starter/internal/domain/ports"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

// corsMiddleware defines CORS specifications
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs each HTTP request with method, URL, and response time
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture the start time
		start := time.Now()

		// Log the incoming request details
		slog.Info("REQUEST",
			"url", r.URL.String(),
			"method", r.Method,
		)

		// Wrap the ResponseWriter to capture the status code
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrappedWriter, r)

		// Log the response details
		duration := time.Since(start)
		slog.Info("RESPONSE",
			"url", r.URL.String(),
			"method", r.Method,
			"status", wrappedWriter.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

// responseWriter is a wrapper around http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

const (
	// authorizationHeaderKey defines the key used to retrieve the authorization header from the HTTP request.
	authorizationHeaderKey = "Authorization"
	// authorizationType specifies the accepted type of authorization.
	authorizationType = "bearer"
	// authorizationPayloadKey defines the key used to store and retrieve the authorization payload from the context.
	authorizationPayloadKey = "authorization_payload"
)

// authMiddleware is a middleware function that validates the authorization token from the incoming HTTP request.
func authMiddleware(tokenService *ports.TokenService) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get(authorizationHeaderKey)
			if len(authorizationHeader) == 0 {
				handleError(w, domain.ErrUnauthorized)
				return
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) != 2 {
				handleError(w, domain.ErrUnauthorized)
				return
			}

			if strings.ToLower(fields[0]) != strings.ToLower(authorizationType) {
				handleError(w, domain.ErrUnauthorized)
				return
			}

			accessToken := fields[1]
			tokenPayload, err := (*tokenService).ValidateAndParseAccessToken(accessToken)
			if err != nil {
				handleError(w, domain.ErrUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), authorizationPayloadKey, tokenPayload)
			r = r.WithContext(ctx)

			f(w, r)
		}
	}
}

// guestMiddleware is a middleware function that allows access to HTTP requests from guests (unauthenticated users).
// It checks for the presence of an authorization token in the request header. If a valid token is found,
// it responds with a forbidden error, preventing authenticated users from accessing guest-only routes.
func guestMiddleware(tokenService *ports.TokenService) Middleware {
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
				handleError(w, domain.ErrForbidden)
				return
			}

			f(w, r)
		}
	}
}
