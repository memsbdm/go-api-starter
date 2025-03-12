package middleware

import (
	"context"
	"encoding/json"
	"go-starter/internal/adapters"
	"go-starter/internal/adapters/ratelimiter"
	"go-starter/internal/adapters/server/helpers"
	"go-starter/internal/adapters/server/responses"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/services"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

// RouteMiddleware is a middleware that applies route-specific middleware functions to the HTTP request pipeline.
type RouteMiddleware struct {
	MailLimiter Middleware
	Auth        Middleware
	Admin       Middleware
}

// NewRouteMiddleware creates a new RouteMiddleware instance.
// It initializes the middleware components with the provided services and adapters.
func NewRouteMiddleware(s *services.Services, a *adapters.Adapters) *RouteMiddleware {
	mailLimiter := ratelimiter.New(a.CacheRepository, "mail")
	mailLimiterMiddleware := RateLimitMiddleware(mailLimiter, RateLimitConfig{
		Limit:  1,
		Window: time.Minute,
	}, a.ErrTrackerAdapter)

	return &RouteMiddleware{
		MailLimiter: mailLimiterMiddleware,
		Auth:        AuthMiddleware(s.TokenService, a.ErrTrackerAdapter),
		Admin:       RoleMiddleware(s.UserService, AuthMiddleware(s.TokenService, a.ErrTrackerAdapter), entities.RoleAdmin),
	}
}

// RateLimitConfig defines configuration for rate limiting
type RateLimitConfig struct {
	Limit  int64
	Window time.Duration
}

// RateLimitMiddleware creates a middleware that limits request rates
func RateLimitMiddleware(limiter *ratelimiter.RateLimiter, config RateLimitConfig, errTracker ports.ErrTrackerAdapter) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identifier := r.RemoteAddr

			// Check if we have an X-Forwarded-For header (case of a proxy)
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				ips := strings.Split(forwardedFor, ",")
				if len(ips) > 0 {
					identifier = strings.TrimSpace(ips[0])
				}
			}

			result, err := limiter.Check(r.Context(), identifier, config.Limit, config.Window)
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

			// Continue to the next handler if rate limit not exceeded
			next.ServeHTTP(w, r)
		})
	}
}

// RoleMiddleware is a middleware function that checks if the user has the required role to access the resource.
func RoleMiddleware(userSvc ports.UserService, authMiddleware Middleware, roleIDs ...entities.RoleID) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		// apply auth middleware first to ensure we have a valid user id in the context
		handlerWithAuth := authMiddleware(func(w http.ResponseWriter, r *http.Request) {
			userID, err := helpers.GetUserIDFromContext(r.Context())
			if err != nil {
				responses.HandleError(w, domain.ErrUnauthorized)
				return
			}

			user, err := userSvc.GetByID(r.Context(), userID)
			if err != nil {
				responses.HandleError(w, domain.ErrUnauthorized)
				return
			}

			hasValidRole := slices.Contains(roleIDs, user.RoleID)

			if !hasValidRole {
				responses.HandleError(w, domain.ErrForbidden)
				return
			}

			f(w, r)
		})

		return handlerWithAuth
	}
}

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
