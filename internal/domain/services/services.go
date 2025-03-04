package services

import (
	"go-starter/config"
	"go-starter/internal/adapters"
	"go-starter/internal/domain/ports"
)

// Services holds all service implementations for the application.
type Services struct {
	CacheService      ports.CacheService
	UserService       ports.UserService
	AuthService       ports.AuthService
	TokenService      ports.TokenService
	MailerService     ports.MailerService
	FileUploadService ports.FileUploadService
}

// New creates and initializes a new Services instance with the provided dependencies.
func New(cfg *config.Container, a *adapters.Adapters) *Services {
	fileUploadSvc := NewFileUploadService(a.FileUploadAdapter)
	cacheSvc := NewCacheService(a.CacheRepository)
	tokenSvc := NewTokenService(cfg.Token, a.TokenRepository, cacheSvc)
	mailerSvc := NewMailerService(cfg, a.MailerAdapter)
	userSvc := NewUserService(cfg.Application, a.UserRepository, cacheSvc, tokenSvc, mailerSvc, fileUploadSvc)
	authSvc := NewAuthService(cfg.Application, userSvc, tokenSvc, mailerSvc)
	return &Services{
		CacheService:      cacheSvc,
		UserService:       userSvc,
		AuthService:       authSvc,
		TokenService:      tokenSvc,
		MailerService:     mailerSvc,
		FileUploadService: fileUploadSvc,
	}
}
