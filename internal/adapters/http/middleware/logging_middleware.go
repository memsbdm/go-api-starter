package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs each HTTP request with method, URL, and response time
func LoggingMiddleware(next http.Handler) http.Handler {
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
