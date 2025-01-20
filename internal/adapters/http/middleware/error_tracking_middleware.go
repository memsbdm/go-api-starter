package middleware

import (
	"bytes"
	"go-starter/internal/domain/ports"
	"io"
	"log/slog"
	"net/http"
)

// ErrTrackingMiddleware creates a middleware that integrates error tracking functionality
// into the HTTP request pipeline. It captures request details and bodies for error monitoring.
func ErrTrackingMiddleware(errTracker ports.ErrorTracker) HandlerMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip body reading for GET/HEAD requests
			if r.Method == http.MethodGet || r.Method == http.MethodHead {
				errTracker.SetRequest(r)
				next.ServeHTTP(w, r)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("failed to read request body",
					"error", err,
					"path", r.URL.Path,
					"method", r.Method,
				)
				errTracker.CaptureException(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Replace the body for downstream handlers
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			// Track request with body
			errTracker.SetRequest(r)
			errTracker.SetBody(body)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
