//go:build !integration

package services_test

import (
	"go-starter/config"
	"go-starter/internal/adapters/storage/postgres/repositories/mocks"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/adapters/token"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/services"
)

type TestBuilder struct {
	TimeGenerator ports.TimeGenerator
	CacheRepo     ports.CacheRepository
	UserRepo      ports.UserRepository
	TokenRepo     ports.TokenRepository
	CacheService  ports.CacheService
	UserService   ports.UserService
	TokenService  ports.TokenService
	AuthService   ports.AuthService
	TokenConfig   *config.Token
	ErrTracker    ports.ErrorTracker
}

func NewTestBuilder() *TestBuilder {
	errTracker := mocks.NewErrorTrackerMock(&config.ErrTracker{})
	timeGenerator := timegen.NewRealTimeGenerator()
	cacheRepo := mocks.NewCacheMock(timeGenerator)
	tokenRepo := token.NewTokenRepository(timeGenerator, errTracker)
	userRepo := mocks.NewUserRepositoryMock()
	tokenConfig := &config.Token{
		AccessTokenDuration:   accessTokenExpirationDuration,
		RefreshTokenDuration:  refreshTokenExpirationDuration,
		AccessTokenSignature:  []byte("access"),
		RefreshTokenSignature: []byte("refresh"),
	}
	return &TestBuilder{
		TimeGenerator: timeGenerator,
		CacheRepo:     cacheRepo,
		UserRepo:      userRepo,
		TokenRepo:     tokenRepo,
		TokenConfig:   tokenConfig,
		ErrTracker:    errTracker,
	}
}

func (tb *TestBuilder) WithTimeGenerator(tg ports.TimeGenerator) *TestBuilder {
	tb.TimeGenerator = tg
	tb.CacheRepo = mocks.NewCacheMock(tg)
	tb.TokenRepo = token.NewTokenRepository(tg, tb.ErrTracker)
	return tb
}

func (tb *TestBuilder) Build() *TestBuilder {
	tb.CacheService = services.NewCacheService(tb.CacheRepo, tb.ErrTracker)
	tb.UserService = services.NewUserService(tb.UserRepo, tb.CacheService, tb.ErrTracker)
	tb.TokenService = services.NewTokenService(tb.TokenConfig, tb.TokenRepo, tb.CacheService, tb.ErrTracker)
	tb.AuthService = services.NewAuthService(tb.UserService, tb.TokenService, tb.ErrTracker)
	return tb
}
