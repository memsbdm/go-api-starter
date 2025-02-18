package handlers

import (
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/services"
)

// Handlers holds all handler implementations for the application.
type Handlers struct {
	HealthHandler *HealthHandler
	AuthHandler   *AuthHandler
	UserHandler   *UserHandler
	MailerHandler *MailerHandler
	ErrTracker    ports.ErrorTracker
}

// New creates and initializes a new Handlers instance with the provided dependencies.
func New(s *services.Services) *Handlers {
	return &Handlers{
		HealthHandler: NewHealthHandler(s.ErrTracker),
		AuthHandler:   NewAuthHandler(s.AuthService, s.ErrTracker),
		UserHandler:   NewUserHandler(s.UserService, s.ErrTracker),
		MailerHandler: NewMailerHandler(s.ErrTracker, s.MailerService),
		ErrTracker:    s.ErrTracker,
	}
}
