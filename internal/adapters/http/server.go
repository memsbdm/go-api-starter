package http

import (
	"fmt"
	"go-starter/config"
	_ "go-starter/docs"
	"go-starter/internal/adapters/http/handlers"
	m "go-starter/internal/adapters/http/middleware"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"log/slog"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Server is a wrapper for HTTP server.
type Server struct {
	*http.Server
	mux        *http.ServeMux
	handlers   *handlers.Handlers
	errTracker ports.ErrTrackerAdapter
}

// New creates and initializes a new HTTP server.
func New(
	httpConfig *config.HTTP,
	handlers *handlers.Handlers,
	tokenSvc ports.TokenService,
	userSvc ports.UserService,
	errTracker ports.ErrTrackerAdapter,
) *Server {
	server := &Server{
		mux:        http.NewServeMux(),
		handlers:   handlers,
		errTracker: errTracker,
	}

	// Configure routes
	server.setupRoutes(tokenSvc, userSvc)

	// Configure server
	server.Server = &http.Server{
		Addr:         fmt.Sprintf(":%d", httpConfig.Port),
		Handler:      server.setupMiddleware(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) setupRoutes(tokenSvc ports.TokenService, userSvc ports.UserService) {
	auth := m.AuthMiddleware(tokenSvc, s.errTracker)
	adminRole := m.RoleMiddleware(userSvc, auth, entities.RoleAdmin)
	// userRole := m.RoleMiddleware(userSvc, auth, entities.RoleUser)
	// adminOrUserRole := m.RoleMiddleware(userSvc, auth, entities.RoleAdmin, entities.RoleUser)

	// Global routes
	s.mux.HandleFunc("GET /v1/swagger/", httpSwagger.WrapHandler)
	s.mux.HandleFunc("GET /v1/health/postgres", s.handlers.HealthHandler.PostgresHealth)
	s.mux.HandleFunc("GET /v1/mailer", m.Chain(s.handlers.MailerHandler.SendEmail, adminRole))

	// Auth routes
	s.mux.HandleFunc("POST /v1/auth/login", s.handlers.AuthHandler.Login)
	s.mux.HandleFunc("POST /v1/auth/register", s.handlers.AuthHandler.Register)
	s.mux.HandleFunc("DELETE /v1/auth/logout", m.Chain(s.handlers.AuthHandler.Logout, auth))
	s.mux.HandleFunc("POST /v1/auth/password-reset", s.handlers.AuthHandler.SendPasswordResetEmail)
	s.mux.HandleFunc("GET /v1/auth/password-reset/{token}", s.handlers.AuthHandler.VerifyPasswordResetToken)
	s.mux.HandleFunc("PATCH /v1/auth/password-reset/{token}", s.handlers.AuthHandler.ResetPassword)

	// User routes
	s.mux.HandleFunc("GET /v1/users/me", m.Chain(s.handlers.UserHandler.Me, auth))
	s.mux.HandleFunc("POST /v1/users/me/avatar", m.Chain(s.handlers.UserHandler.UploadAvatar, auth))
	s.mux.HandleFunc("PATCH /v1/users/me/password", m.Chain(s.handlers.UserHandler.UpdatePassword, auth))
	s.mux.HandleFunc("GET /v1/users/me/verify-email/{token}", s.handlers.UserHandler.VerifyEmail)
	s.mux.HandleFunc("POST /v1/users/me/verify-email/resend", m.Chain(s.handlers.UserHandler.ResendEmailVerification, auth))
	s.mux.HandleFunc("GET /v1/users/{uuid}", s.handlers.UserHandler.GetByID)
}

func (s *Server) setupMiddleware() http.Handler {
	routerMiddleware := []m.HandlerMiddleware{
		m.ErrTrackingMiddleware(s.errTracker),
		m.LoggingMiddleware(),
		m.SecurityHeadersMiddleware(),
		m.CorsMiddleware(),
	}

	router := m.ChainHandlerFunc(s.mux, routerMiddleware...)
	return s.errTracker.Handle(router)
}

// Serve starts the HTTP server and listens for incoming requests.
func (s *Server) Serve() error {
	slog.Info("Starting HTTP server")
	return s.ListenAndServe()
}
