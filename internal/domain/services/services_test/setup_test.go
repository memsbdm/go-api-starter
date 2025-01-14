// //go:build !integration

package services_test

import (
	"go-starter/config"
	"go-starter/internal/adapters/storage/postgres/repositories/mocks"
	"go-starter/internal/adapters/storage/redis"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/adapters/token"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/services"
	"time"
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
}

func NewTestBuilder() *TestBuilder {
	timeGenerator := timegen.NewRealTimeGenerator()
	cacheRepo := redis.NewMock(timeGenerator)
	tokenRepo := token.NewTokenRepository(timeGenerator)
	userRepo := mocks.MockUserRepository()
	tokenConfig := &config.Token{
		AccessTokenDuration:  10 * time.Minute,
		RefreshTokenDuration: 1 * time.Hour,
		JWTSecret:            []byte("secret"),
	}
	return &TestBuilder{
		TimeGenerator: timeGenerator,
		CacheRepo:     cacheRepo,
		UserRepo:      userRepo,
		TokenRepo:     tokenRepo,
		TokenConfig:   tokenConfig,
	}
}

func (tb *TestBuilder) WithTimeGenerator(tg ports.TimeGenerator) *TestBuilder {
	tb.TimeGenerator = tg
	tb.CacheRepo = redis.NewMock(tg)
	tb.TokenRepo = token.NewTokenRepository(tg)
	return tb
}

func (tb *TestBuilder) Build() *TestBuilder {
	tb.CacheService = services.NewCacheService(tb.CacheRepo)
	tb.UserService = services.NewUserService(tb.UserRepo, tb.CacheService)
	tb.TokenService = services.NewTokenService(tb.TokenConfig, tb.TokenRepo, tb.CacheService)
	tb.AuthService = services.NewAuthService(tb.UserService, tb.TokenService)
	return tb
}
