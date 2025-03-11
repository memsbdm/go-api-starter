//go:build !integration

package services_test

import (
	"go-starter/config"
	"go-starter/internal/adapters/errtracker"
	"go-starter/internal/adapters/mailer"
	"go-starter/internal/adapters/storage/cache"
	"go-starter/internal/adapters/storage/fileupload"
	"go-starter/internal/adapters/storage/postgres/repositories"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/adapters/token"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/services"
)

// debugEmail is the address used for mail sent in dev mode.
// All emails in non-production environments will be redirected to this address.
const debugEmail = "debug@example.com"

type TestBuilder struct {
	TimeGenerator     ports.TimeGenerator
	CacheRepo         ports.CacheRepository
	UserRepo          ports.UserRepository
	TokenProvider     ports.TokenProvider
	CacheService      ports.CacheService
	UserService       ports.UserService
	TokenService      ports.TokenService
	AuthService       ports.AuthService
	Config            *config.Container
	ErrTrackerAdapter ports.ErrTrackerAdapter
	MailerService     ports.MailerService
	MailerAdapter     ports.MailerAdapter
	FileUploadAdapter ports.FileUploadAdapter
	FileUploadService ports.FileUploadService
}

func NewTestBuilder() *TestBuilder {
	fileUploadAdapter := fileupload.NewFileUploadAdapterMock()
	mailerAdapter := mailer.NewMailerAdapterMock()
	errTrackerAdapter := errtracker.NewErrTrackerAdapterMock()
	timeGenerator := timegen.NewTimeGenerator()
	cacheRepo := cache.NewCacheRepositoryMock(timeGenerator)
	tokenProvider := token.NewTokenProvider(timeGenerator, errTrackerAdapter)
	userRepo := repositories.NewUserRepositoryMock()

	cfg := setConfig()

	return &TestBuilder{
		TimeGenerator:     timeGenerator,
		CacheRepo:         cacheRepo,
		UserRepo:          userRepo,
		TokenProvider:     tokenProvider,
		Config:            cfg,
		ErrTrackerAdapter: errTrackerAdapter,
		MailerAdapter:     mailerAdapter,
		FileUploadAdapter: fileUploadAdapter,
	}
}

func (tb *TestBuilder) WithTimeGenerator(tg ports.TimeGenerator) *TestBuilder {
	tb.TimeGenerator = tg
	tb.CacheRepo = cache.NewCacheRepositoryMock(tg)
	tb.TokenProvider = token.NewTokenProvider(tg, tb.ErrTrackerAdapter)
	return tb
}

func (tb *TestBuilder) Build() *TestBuilder {
	tb.FileUploadService = services.NewFileUploadService(tb.FileUploadAdapter)
	tb.MailerService = services.NewMailerService(tb.Config, tb.MailerAdapter)
	tb.CacheService = services.NewCacheService(tb.CacheRepo)
	tb.TokenService = services.NewTokenService(tb.Config.Token, tb.TokenProvider, tb.CacheService)
	tb.UserService = services.NewUserService(tb.Config, tb.UserRepo, tb.CacheService, tb.TokenService, tb.MailerService, tb.FileUploadService)
	tb.AuthService = services.NewAuthService(tb.Config, tb.UserService, tb.TokenService, tb.MailerService)
	return tb
}

func (tb *TestBuilder) SetEnvToProduction() *TestBuilder {
	tb.Config.Application.Env = config.EnvProduction
	return tb
}

func setConfig() *config.Container {
	appConfig := &config.App{
		Env: config.EnvDevelopment,
	}

	tokenConfig := &config.Token{
		AccessTokenDuration:            accessTokenExpirationDuration,
		EmailVerificationTokenDuration: emailVerificationTokenExpirationDuration,
		PasswordResetTokenDuration:     passwordResetTokenExpirationDuration,
	}

	mailerConfig := &config.Mailer{
		DebugTo: debugEmail,
	}

	return &config.Container{
		Application: appConfig,
		Token:       tokenConfig,
		Mailer:      mailerConfig,
	}
}
