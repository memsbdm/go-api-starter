package http

import (
	"errors"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"go-starter/config"
	_ "go-starter/docs"
	"log/slog"
	"net/http"
	"time"
)

// Server is a wrapper for HTTP server
type Server struct {
	*http.Server
}

// New creates a new HTTP server
func New(config *config.HTTP, userHandler UserHandler) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /v1/users/{id}", userHandler.GetByID)
	mux.HandleFunc("POST /v1/users", userHandler.Register)

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
