package http

import (
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"go-starter/config"
	_ "go-starter/docs"
	"go-starter/internal/adapters/http/handlers"
	m "go-starter/internal/adapters/http/middleware"
	"go-starter/internal/domain/ports"
	"log/slog"
	"net/http"
	"time"
)

// Server is a wrapper for HTTP server.
type Server struct {
	*http.Server
}

// New creates and initializes a new HTTP server.
func New(
	httpConfig *config.HTTP,
	handlers *handlers.Handlers,
	tokenService ports.TokenService,
	errTrackerAdapter ports.ErrTrackerAdapter,
) *Server {
	auth := func() m.Middleware {
		return m.AuthMiddleware(tokenService, errTrackerAdapter)
	}
	guest := func() m.Middleware {
		return m.GuestMiddleware(tokenService)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /v1/health", handlers.HealthHandler.Health)

	// Mailer
	mux.HandleFunc("GET /v1/mailer", handlers.MailerHandler.SendEmail)

	// Auth
	mux.HandleFunc("POST /v1/auth/login", m.Chain(handlers.AuthHandler.Login, guest()))
	mux.HandleFunc("POST /v1/auth/register", m.Chain(handlers.AuthHandler.Register, guest()))
	mux.HandleFunc("POST /v1/auth/refresh", handlers.AuthHandler.Refresh)
	mux.HandleFunc("DELETE /v1/auth/logout", m.Chain(handlers.AuthHandler.Logout, auth()))

	// Users
	mux.HandleFunc("GET /v1/users/me", m.Chain(handlers.UserHandler.Me, auth()))
	mux.HandleFunc("GET /v1/users/{uuid}", handlers.UserHandler.GetByID)
	mux.HandleFunc("PATCH /v1/users/password", m.Chain(handlers.UserHandler.UpdatePassword, auth()))

	routerMiddleware := []m.HandlerMiddleware{
		m.ErrTrackingMiddleware(errTrackerAdapter),
		m.LoggingMiddleware(),
		m.SecurityHeadersMiddleware(),
		m.CorsMiddleware(),
	}

	router := m.ChainHandlerFunc(mux, routerMiddleware...)
	handler := errTrackerAdapter.Handle(router)

	return &Server{
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%d", httpConfig.Port),
			Handler:      handler,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}
}

// Serve starts the HTTP server and listens for incoming requests.
func (s *Server) Serve() error {
	slog.Info("Starting HTTP server")
	return s.ListenAndServe()
}
