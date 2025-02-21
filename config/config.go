package config

import (
	"go-starter/pkg/env"
	"time"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

type (
	// Container contains environment variables for the application, database, http server, ...
	Container struct {
		Application *App
		DB          *DB
		HTTP        *HTTP
		Redis       *Redis
		Token       *Token
		ErrTracker  *ErrTracker
		Mailer      *Mailer
	}

	// App contains all the environment variables for the application.
	App struct {
		Env string
	}

	// DB contains all the environment variables for the database.
	DB struct {
		Addr         string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}

	// HTTP contains all the environment variables for the http server.
	HTTP struct {
		Port int
	}

	// Redis contains all the environment variables for the cache service.
	Redis struct {
		Addr     string
		Password string
	}

	// Token contains all the environment variables for the token service.
	Token struct {
		AccessTokenSignature  []byte
		RefreshTokenSignature []byte
		AccessTokenDuration   time.Duration
		RefreshTokenDuration  time.Duration
	}

	// ErrTracker contains all the environment variables for the error tracking.
	ErrTracker struct {
		DSN              string
		TracesSampleRate float64
	}

	// Mailer contains all the environment variables for the mailer.
	Mailer struct {
		Host                string
		Port                int
		Username            string
		Password            string
		From                string
		DebugTo             string
		MaxRetries          int
		RetryDelayInSeconds int
	}
)

// New creates a new container instance.
func New() *Container {
	app := &App{
		Env: env.GetString("ENVIRONMENT"),
	}

	db := &DB{
		Addr:         env.GetString("DB_ADDR"),
		MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS"),
		MaxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME"),
	}

	http := &HTTP{
		Port: env.GetInt("HTTP_PORT"),
	}

	redis := &Redis{
		Addr:     env.GetString("REDIS_ADDR"),
		Password: env.GetString("REDIS_PASSWORD"),
	}

	token := &Token{
		AccessTokenSignature:  []byte(env.GetString("ACCESS_TOKEN_SIGNATURE")),
		RefreshTokenSignature: []byte(env.GetString("REFRESH_TOKEN_SIGNATURE")),
		AccessTokenDuration:   env.GetDuration("ACCESS_TOKEN_DURATION"),
		RefreshTokenDuration:  env.GetDuration("REFRESH_TOKEN_DURATION"),
	}

	errTracker := &ErrTracker{
		DSN:              env.GetString("SENTRY_DSN"),
		TracesSampleRate: env.GetFloat64("SENTRY_TRACES_SAMPLE_RATE"),
	}

	mailer := &Mailer{
		Host:                env.GetString("MAILER_HOST"),
		Port:                env.GetInt("MAILER_PORT"),
		Username:            env.GetString("MAILER_USERNAME"),
		Password:            env.GetString("MAILER_PASSWORD"),
		From:                env.GetString("MAILER_FROM"),
		DebugTo:             env.GetString("MAILER_DEBUG_TO"),
		MaxRetries:          env.GetInt("MAILER_MAX_RETRIES"),
		RetryDelayInSeconds: env.GetInt("MAILER_RETRIES_DELAY_IN_SECONDS"),
	}

	return &Container{
		Application: app,
		DB:          db,
		HTTP:        http,
		Redis:       redis,
		Token:       token,
		ErrTracker:  errTracker,
		Mailer:      mailer,
	}
}
