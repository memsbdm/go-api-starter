package server

import (
	"go-starter/internal/adapters"
	"go-starter/internal/adapters/server/handlers"
	m "go-starter/internal/adapters/server/middleware"
	"go-starter/internal/domain/services"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(h *handlers.Handlers, s *services.Services, a *adapters.Adapters) http.Handler {
	mux := http.NewServeMux()

	// Global middleware
	gm := m.NewGlobalMiddleware(a.ErrTrackerAdapter, a.CacheRepository)
	handler := m.ChainHandlerFunc(mux,
		gm.ErrTracking,
		gm.Logging,
		gm.RateLimiter,
		gm.Security,
		gm.Cors,
	)
	handler = a.ErrTrackerAdapter.Handle(handler)

	// Route middleware
	rm := m.NewRouteMiddleware(s, a)

	// Global routes
	mux.HandleFunc("GET /v1/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /v1/health/postgres", m.Chain(h.HealthHandler.PostgresHealth))
	mux.HandleFunc("GET /v1/mailer", m.Chain(h.MailerHandler.SendEmail, rm.Admin))

	// Auth routes
	mux.HandleFunc("POST /v1/auth/login", h.AuthHandler.Login)
	mux.HandleFunc("POST /v1/auth/register", h.AuthHandler.Register)
	mux.HandleFunc("DELETE /v1/auth/logout", m.Chain(h.AuthHandler.Logout, rm.Auth))
	mux.HandleFunc("POST /v1/auth/password-reset", m.Chain(h.AuthHandler.SendPasswordResetEmail, rm.MailLimiter))
	mux.HandleFunc("GET /v1/auth/password-reset/{token}", h.AuthHandler.VerifyPasswordResetToken)
	mux.HandleFunc("PATCH /v1/auth/password-reset/{token}", h.AuthHandler.ResetPassword)

	// User routes
	mux.HandleFunc("GET /v1/users/me", m.Chain(h.UserHandler.Me, rm.Auth))
	mux.HandleFunc("POST /v1/users/me/avatar", m.Chain(h.UserHandler.UploadAvatar, rm.Auth))
	mux.HandleFunc("DELETE /v1/users/me/avatar", m.Chain(h.UserHandler.DeleteAvatar, rm.Auth))
	mux.HandleFunc("PATCH /v1/users/me/password", m.Chain(h.UserHandler.UpdatePassword, rm.Auth))
	mux.HandleFunc("GET /v1/users/me/verify-email/{token}", h.UserHandler.VerifyEmail)
	mux.HandleFunc("POST /v1/users/me/verify-email/resend", m.Chain(h.UserHandler.ResendEmailVerification, rm.Auth, rm.MailLimiter))
	mux.HandleFunc("GET /v1/users/{uuid}", h.UserHandler.GetByID)

	return handler
}
