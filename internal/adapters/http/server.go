package http

import (
	"errors"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"go-starter/config"
	_ "go-starter/docs"
	"go-starter/internal/adapters/http/handlers"
	"go-starter/internal/adapters/http/middleware"
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
func New(config *config.HTTP, healthHandler handlers.HealthHandler, authHandler handlers.AuthHandler, userHandler handlers.UserHandler, tokenService ports.TokenService) *Server {
	auth := func() middleware.Middleware {
		return middleware.AuthMiddleware(&tokenService)
	}
	guest := func() middleware.Middleware {
		return middleware.GuestMiddleware(&tokenService)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /v1/health", healthHandler.Health)

	// Auth
	mux.HandleFunc("POST /v1/auth/login", middleware.Chain(authHandler.Login, guest()))
	mux.HandleFunc("POST /v1/auth/register", middleware.Chain(authHandler.Register, guest()))
	mux.HandleFunc("POST /v1/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("DELETE /v1/auth/logout", middleware.Chain(authHandler.Logout, auth()))

	// Users
	mux.HandleFunc("GET /v1/users/me", middleware.Chain(userHandler.Me, auth()))
	mux.HandleFunc("GET /v1/users/{uuid}", userHandler.GetByID)
	mux.HandleFunc("PATCH /v1/users/password", middleware.Chain(userHandler.UpdatePassword, auth()))

	handler := middleware.LoggingMiddleware(middleware.CorsMiddleware(mux))
	return &Server{
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%d", config.Port),
			Handler:      handler,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}
}

// Serve starts the HTTP server and listens for incoming requests.
func (s *Server) Serve() {
	slog.Info("Starting HTTP server")
	err := s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
