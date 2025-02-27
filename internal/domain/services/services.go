package services

import (
	"go-starter/config"
	"go-starter/internal/adapters"
	"go-starter/internal/domain/ports"
)

// Services holds all service implementations for the application.
type Services struct {
	CacheService  ports.CacheService
	UserService   ports.UserService
	AuthService   ports.AuthService
	TokenService  ports.TokenService
	MailerService ports.MailerService
}

// New creates and initializes a new Services instance with the provided dependencies.
func New(cfg *config.Container, a *adapters.Adapters) *Services {
	cacheSvc := NewCacheService(a.CacheRepository)
	tokenSvc := NewTokenService(cfg.Token, a.TokenRepository, cacheSvc)
	userSvc := NewUserService(a.UserRepository, cacheSvc, tokenSvc)
	mailerSvc := NewMailerService(cfg, a.MailerAdapter)
	authSvc := NewAuthService(cfg.Application, userSvc, tokenSvc, mailerSvc)
	return &Services{
		CacheService:  cacheSvc,
		UserService:   userSvc,
		AuthService:   authSvc,
		TokenService:  tokenSvc,
		MailerService: mailerSvc,
	}
}
