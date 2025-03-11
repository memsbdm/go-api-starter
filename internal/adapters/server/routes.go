package server

import (
	"go-starter/internal/adapters/server/handlers"
	m "go-starter/internal/adapters/server/middleware"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/services"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(h *handlers.Handlers, s *services.Services, errTracker ports.ErrTrackerAdapter) http.Handler {
	mux := http.NewServeMux()

	// Global middleware
	routerMiddleware := []m.HandlerMiddleware{
		m.ErrTrackingMiddleware(errTracker),
		m.LoggingMiddleware(),
		m.SecurityHeadersMiddleware(),
		m.CorsMiddleware(),
	}
	handler := m.ChainHandlerFunc(mux, routerMiddleware...)
	handler = errTracker.Handle(handler)

	// Routes middleware
	auth := m.AuthMiddleware(s.TokenService, errTracker)
	adminRole := m.RoleMiddleware(s.UserService, auth, entities.RoleAdmin)

	// Global routes
	mux.HandleFunc("GET /v1/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /v1/health/postgres", h.HealthHandler.PostgresHealth)
	mux.HandleFunc("GET /v1/mailer", m.Chain(h.MailerHandler.SendEmail, adminRole))

	// Auth routes
	mux.HandleFunc("POST /v1/auth/login", h.AuthHandler.Login)
	mux.HandleFunc("POST /v1/auth/register", h.AuthHandler.Register)
	mux.HandleFunc("DELETE /v1/auth/logout", m.Chain(h.AuthHandler.Logout, auth))
	mux.HandleFunc("POST /v1/auth/password-reset", h.AuthHandler.SendPasswordResetEmail)
	mux.HandleFunc("GET /v1/auth/password-reset/{token}", h.AuthHandler.VerifyPasswordResetToken)
	mux.HandleFunc("PATCH /v1/auth/password-reset/{token}", h.AuthHandler.ResetPassword)

	// User routes
	mux.HandleFunc("GET /v1/users/me", m.Chain(h.UserHandler.Me, auth))
	mux.HandleFunc("POST /v1/users/me/avatar", m.Chain(h.UserHandler.UploadAvatar, auth))
	mux.HandleFunc("DELETE /v1/users/me/avatar", m.Chain(h.UserHandler.DeleteAvatar, auth))
	mux.HandleFunc("PATCH /v1/users/me/password", m.Chain(h.UserHandler.UpdatePassword, auth))
	mux.HandleFunc("GET /v1/users/me/verify-email/{token}", h.UserHandler.VerifyEmail)
	mux.HandleFunc("POST /v1/users/me/verify-email/resend", m.Chain(h.UserHandler.ResendEmailVerification, auth))
	mux.HandleFunc("GET /v1/users/{uuid}", h.UserHandler.GetByID)

	return handler
}
