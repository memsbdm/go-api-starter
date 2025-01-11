package http

import (
	"errors"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"go-starter/config"
	_ "go-starter/docs"
	"go-starter/internal/domain/ports"
	"log/slog"
	"net/http"
	"time"
)

// Server is a wrapper for HTTP server
type Server struct {
	*http.Server
}

// New creates a new HTTP server
func New(config *config.HTTP, healthHandler HealthHandler, authHandler AuthHandler, userHandler UserHandler, tokenService ports.TokenService) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /v1/health", healthHandler.Health)
	mux.HandleFunc("POST /v1/auth/login", authHandler.Login)
	mux.Handle("GET /v1/users/{uuid}", Chain(userHandler.GetByID, authMiddleware(&tokenService)))
	mux.Handle("POST /v1/users", Chain(userHandler.Register, guestMiddleware(&tokenService)))

	handler := loggingMiddleware(corsMiddleware(mux))
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

// Serve starts the HTTP server
func (s *Server) Serve() {
	slog.Info("Starting HTTP server")
	err := s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
