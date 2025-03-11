package server

import (
	"fmt"
	"go-starter/config"
	_ "go-starter/docs"
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
	handler http.Handler,
) *Server {
	server := &Server{}

	// Configure server
	server.Server = &http.Server{
		Addr:         fmt.Sprintf(":%d", httpConfig.Port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

// Serve starts the HTTP server and listens for incoming requests.
func (s *Server) Serve() error {
	slog.Info("starting HTTP server")
	return s.ListenAndServe()
}
