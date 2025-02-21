package handlers

import (
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/services"
)

// Handlers holds all handler implementations for the application.
type Handlers struct {
	HealthHandler     *HealthHandler
	AuthHandler       *AuthHandler
	UserHandler       *UserHandler
	MailerHandler     *MailerHandler
	ErrTrackerAdapter ports.ErrTrackerAdapter
}

// New creates and initializes a new Handlers instance with the provided dependencies.
func New(s *services.Services) *Handlers {
	return &Handlers{
		HealthHandler:     NewHealthHandler(s.ErrTrackerAdapter),
		AuthHandler:       NewAuthHandler(s.AuthService, s.ErrTrackerAdapter),
		UserHandler:       NewUserHandler(s.UserService, s.ErrTrackerAdapter),
		MailerHandler:     NewMailerHandler(s.ErrTrackerAdapter, s.MailerService),
		ErrTrackerAdapter: s.ErrTrackerAdapter,
	}
}
