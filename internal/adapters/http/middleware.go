package http

import (
	"context"
	"go-starter/internal/domain/ports"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

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
	// authorizationHeaderKey is the key for authorization header in the request
	authorizationHeaderKey = "Authorization"
	// authorizationType is the accepted authorization type
	authorizationType = "bearer"
	// authorizationPayloadKey is the key for authorization payload in the context
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(authService *ports.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if strings.ToLower(fields[0]) != strings.ToLower(authorizationType) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken := fields[1]
		tokenPayload, err := (*authService).ValidateToken(accessToken)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authorizationPayloadKey, tokenPayload)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
