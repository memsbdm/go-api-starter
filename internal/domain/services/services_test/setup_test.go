//go:build !integration

package services_test

import (
	"go-starter/config"
	"go-starter/internal/adapters/mocks"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/adapters/token"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/services"
)

// debugEmail is the address used for mail sent in dev mode.
// All emails in non-production environments will be redirected to this address.
const debugEmail = "debug@example.com"

type TestBuilder struct {
	TimeGenerator ports.TimeGenerator
	CacheRepo     ports.CacheRepository
	UserRepo      ports.UserRepository
	TokenRepo     ports.TokenRepository
	CacheService  ports.CacheService
	UserService   ports.UserService
	TokenService  ports.TokenService
	AuthService   ports.AuthService
	Config        *config.Container
	ErrTracker    ports.ErrorTracker
	MailerService ports.MailerService
	MailerRepo    ports.MailerRepository
}

func NewTestBuilder() *TestBuilder {
	mailerRepo := mocks.NewMailerRepositoryMock()
	errTracker := mocks.NewErrorTrackerMock(&config.ErrTracker{})
	timeGenerator := timegen.NewRealTimeGenerator()
	cacheRepo := mocks.NewCacheMock(timeGenerator)
	tokenRepo := token.NewTokenRepository(timeGenerator, errTracker)
	userRepo := mocks.NewUserRepositoryMock()

	cfg := setConfig()

	return &TestBuilder{
		TimeGenerator: timeGenerator,
		CacheRepo:     cacheRepo,
		UserRepo:      userRepo,
		TokenRepo:     tokenRepo,
		Config:        cfg,
		ErrTracker:    errTracker,
		MailerRepo:    mailerRepo,
	}
}

func (tb *TestBuilder) WithTimeGenerator(tg ports.TimeGenerator) *TestBuilder {
	tb.TimeGenerator = tg
	tb.CacheRepo = mocks.NewCacheMock(tg)
	tb.TokenRepo = token.NewTokenRepository(tg, tb.ErrTracker)
	return tb
}

func (tb *TestBuilder) Build() *TestBuilder {
	tb.MailerService = services.NewMailerService(tb.Config, tb.MailerRepo, tb.ErrTracker)
	tb.CacheService = services.NewCacheService(tb.CacheRepo, tb.ErrTracker)
	tb.UserService = services.NewUserService(tb.UserRepo, tb.CacheService, tb.ErrTracker)
	tb.TokenService = services.NewTokenService(tb.Config.Token, tb.TokenRepo, tb.CacheService, tb.ErrTracker)
	tb.AuthService = services.NewAuthService(tb.UserService, tb.TokenService, tb.ErrTracker)
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
		AccessTokenDuration:   accessTokenExpirationDuration,
		RefreshTokenDuration:  refreshTokenExpirationDuration,
		AccessTokenSignature:  []byte("access"),
		RefreshTokenSignature: []byte("refresh"),
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
