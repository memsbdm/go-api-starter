package http

import (
	"fmt"
	"go-starter/config"
	_ "go-starter/docs"
	"go-starter/internal/adapters/http/handlers"
	m "go-starter/internal/adapters/http/middleware"
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
	errTracker ports.ErrTrackerAdapter,
) *Server {
	server := &Server{
		mux:        http.NewServeMux(),
		handlers:   handlers,
		errTracker: errTracker,
	}

	// Configure routes
	server.setupRoutes(tokenSvc)

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

func (s *Server) setupRoutes(tokenSvc ports.TokenService) {
	auth := m.AuthMiddleware(tokenSvc, s.errTracker)
	guest := m.GuestMiddleware(tokenSvc)

	// Global routes
	s.mux.HandleFunc("GET /v1/swagger/", httpSwagger.WrapHandler)
	s.mux.HandleFunc("GET /v1/health", s.handlers.HealthHandler.Health)
	s.mux.HandleFunc("GET /v1/mailer", s.handlers.MailerHandler.SendEmail)

	// Auth routes
	s.mux.HandleFunc("POST /v1/auth/login", m.Chain(s.handlers.AuthHandler.Login, guest))
	s.mux.HandleFunc("POST /v1/auth/register", m.Chain(s.handlers.AuthHandler.Register, guest))
	s.mux.HandleFunc("POST /v1/auth/refresh", s.handlers.AuthHandler.Refresh)
	s.mux.HandleFunc("DELETE /v1/auth/logout", m.Chain(s.handlers.AuthHandler.Logout, auth))

	// User routes
	s.mux.HandleFunc("GET /v1/users/me", m.Chain(s.handlers.UserHandler.Me, auth))
	s.mux.HandleFunc("GET /v1/users/{uuid}", s.handlers.UserHandler.GetByID)
	s.mux.HandleFunc("PATCH /v1/users/password", m.Chain(s.handlers.UserHandler.UpdatePassword, auth))
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
