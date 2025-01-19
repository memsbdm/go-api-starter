package handlers

import (
	"go-starter/internal/domain/services"
)

// TODO
type Handlers struct {
	HealthHandler *HealthHandler
	AuthHandler   *AuthHandler
	UserHandler   *UserHandler
}

func New(s *services.Services) *Handlers {
	return &Handlers{
		HealthHandler: NewHealthHandler(),
		AuthHandler:   NewAuthHandler(s.AuthService),
		UserHandler:   NewUserHandler(s.UserService),
	}
}
