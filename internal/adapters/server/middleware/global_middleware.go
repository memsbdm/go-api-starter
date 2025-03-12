package middleware

import (
	"bytes"
	"encoding/json"
	"go-starter/internal/adapters/ratelimiter"
	"go-starter/internal/domain/ports"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GlobalMiddleware is a middleware that applies global middleware functions to the HTTP request pipeline.
type GlobalMiddleware struct {
	ErrTracking HandlerMiddleware
	Logging     HandlerMiddleware
	Security    HandlerMiddleware
	Cors        HandlerMiddleware
	RateLimiter HandlerMiddleware
}

// NewGlobalMiddleware creates a new GlobalMiddleware instance.
// It initializes the middleware components with the provided error tracker and cache repository.
func NewGlobalMiddleware(errTracker ports.ErrTrackerAdapter, cache ports.CacheRepository) *GlobalMiddleware {
	globalLimiter := ratelimiter.New(cache, "global")

	return &GlobalMiddleware{
		ErrTracking: ErrTrackingMiddleware(errTracker),
		Logging:     LoggingMiddleware(),
		Security:    SecurityHeadersMiddleware(),
		Cors:        CorsMiddleware(),
		RateLimiter: GlobalRateLimitMiddleware(globalLimiter, errTracker),
	}
}

// GlobalRateLimitMiddleware creates a middleware that globally limits the request rate.
func GlobalRateLimitMiddleware(limiter *ratelimiter.RateLimiter, errTracker ports.ErrTrackerAdapter) HandlerMiddleware {
	// Configuration du rate limiter global
	const (
		globalLimit  = 200
		globalWindow = time.Minute
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use the IP address as the identifier
			identifier := r.RemoteAddr

			// Check if we have an X-Forwarded-For header (case of a proxy)
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				ips := strings.Split(forwardedFor, ",")
				if len(ips) > 0 {
					identifier = strings.TrimSpace(ips[0])
				}
			}

			result, err := limiter.Check(r.Context(), identifier, globalLimit, globalWindow)
			if err != nil {
				errTracker.CaptureException(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(result.Limit, 10))
			w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(result.Limit-result.Current, 10))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(int64(result.ResetAfter.Seconds()), 10))

			if !result.Allowed {
				w.Header().Set("Retry-After", strconv.FormatInt(int64(result.ResetAfter.Seconds()), 10))
				w.WriteHeader(http.StatusTooManyRequests)

				response := map[string]interface{}{
					"error":       "Rate limit exceeded",
					"retry_after": result.ResetAfter.Seconds(),
				}

				json.NewEncoder(w).Encode(response)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ErrTrackingMiddleware creates a middleware that integrates error tracking functionality
// into the HTTP request pipeline. It captures request details and bodies for error monitoring.
func ErrTrackingMiddleware(errTracker ports.ErrTrackerAdapter) HandlerMiddleware {
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

// CorsMiddleware defines CORS specifications.
func CorsMiddleware() HandlerMiddleware {
	return func(next http.Handler) http.Handler {
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
}

// SecurityHeadersMiddleware sets security header.
func SecurityHeadersMiddleware() HandlerMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
			w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
			w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

			// Proceed with the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware logs each HTTP request with method, URL, and response time
func LoggingMiddleware() HandlerMiddleware {
	return func(next http.Handler) http.Handler {
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
