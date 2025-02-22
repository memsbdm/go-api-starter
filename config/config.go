package config

import (
	"fmt"
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
		MaxOpenConns: env.GetOptionalInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns: env.GetOptionalInt("DB_MAX_IDLE_CONNS"),
		MaxIdleTime:  env.GetOptionalDuration("DB_MAX_IDLE_TIME"),
	}

	http := &HTTP{
		Port: env.GetOptionalInt("HTTP_PORT"),
	}

	redis := &Redis{
		Addr:     env.GetString("REDIS_ADDR"),
		Password: env.GetOptionalString("REDIS_PASSWORD"),
	}

	token := &Token{
		AccessTokenSignature:  []byte(env.GetString("ACCESS_TOKEN_SIGNATURE")),
		RefreshTokenSignature: []byte(env.GetString("REFRESH_TOKEN_SIGNATURE")),
		AccessTokenDuration:   env.GetOptionalDuration("ACCESS_TOKEN_DURATION"),
		RefreshTokenDuration:  env.GetOptionalDuration("REFRESH_TOKEN_DURATION"),
	}

	errTracker := &ErrTracker{
		DSN:              env.GetOptionalString("SENTRY_DSN"),
		TracesSampleRate: env.GetOptionalFloat64("SENTRY_TRACES_SAMPLE_RATE"),
	}

	mailer := &Mailer{
		Host:                env.GetString("MAILER_HOST"),
		Port:                env.GetInt("MAILER_PORT"),
		Username:            env.GetString("MAILER_USERNAME"),
		Password:            env.GetString("MAILER_PASSWORD"),
		From:                env.GetString("MAILER_FROM"),
		DebugTo:             env.GetString("MAILER_DEBUG_TO"),
		MaxRetries:          env.GetOptionalInt("MAILER_MAX_RETRIES"),
		RetryDelayInSeconds: env.GetOptionalInt("MAILER_RETRIES_DELAY_IN_SECONDS"),
	}

	c := &Container{
		Application: app,
		DB:          db,
		HTTP:        http,
		Redis:       redis,
		Token:       token,
		ErrTracker:  errTracker,
		Mailer:      mailer,
	}

	c.setDefaultValues()
	err := c.validate()
	if err != nil {
		panic(err)
	}

	return c
}

// setDefaultValues sets the default values for the container if they are not set.
func (c *Container) setDefaultValues() {
	// DB
	if c.DB.MaxOpenConns == 0 {
		c.DB.MaxOpenConns = 30
	}
	if c.DB.MaxIdleConns == 0 {
		c.DB.MaxIdleConns = 30
	}
	if c.DB.MaxIdleTime == 0 {
		c.DB.MaxIdleTime = 15 * time.Minute
	}

	// HTTP
	if c.HTTP.Port == 0 {
		c.HTTP.Port = 8080
	}

	// Token
	if c.Token.AccessTokenDuration == 0 {
		c.Token.AccessTokenDuration = 15 * time.Minute
	}
	if c.Token.RefreshTokenDuration == 0 {
		c.Token.RefreshTokenDuration = 1 * time.Hour
	}

	// ErrTracker
	if c.ErrTracker.TracesSampleRate == 0 {
		c.ErrTracker.TracesSampleRate = 1.0
	}

	// Mailer
	if c.Mailer.MaxRetries == 0 {
		c.Mailer.MaxRetries = 3
	}
	if c.Mailer.RetryDelayInSeconds == 0 {
		c.Mailer.RetryDelayInSeconds = 10
	}
}

// validate validates the container.
func (c *Container) validate() error {
	// Application
	if c.Application.Env != EnvDevelopment && c.Application.Env != EnvProduction {
		return fmt.Errorf("invalid environment variable: %s", "ENVIRONMENT")
	}

	// DB
	if c.DB.MaxIdleTime < 0 {
		return fmt.Errorf("invalid environment variable: %s", "DB_MAX_IDLE_TIME")
	}

	if c.DB.MaxIdleConns < 0 {
		return fmt.Errorf("invalid environment variable: %s", "DB_MAX_IDLE_CONNS")
	}

	if c.DB.MaxOpenConns < 0 {
		return fmt.Errorf("invalid environment variable: %s", "DB_MAX_OPEN_CONNS")
	}

	if c.DB.MaxOpenConns < c.DB.MaxIdleConns {
		return fmt.Errorf("invalid environment variables: %s should be greater or equal to %s", "DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS")
	}

	// HTTP
	if c.HTTP.Port <= 0 {
		return fmt.Errorf("invalid environment variable: %s", "HTTP_PORT")
	}

	// Token
	if c.Token.AccessTokenDuration < 0 {
		return fmt.Errorf("invalid environment variable: %s", "ACCESS_TOKEN_DURATION")
	}

	if c.Token.RefreshTokenDuration < 0 {
		return fmt.Errorf("invalid environment variable: %s", "REFRESH_TOKEN_DURATION")
	}

	if c.Token.AccessTokenDuration > c.Token.RefreshTokenDuration {
		return fmt.Errorf("invalid environment variables: %s should be less than or equal to %s", "ACCESS_TOKEN_DURATION", "REFRESH_TOKEN_DURATION")
	}

	// ErrTracker
	if c.ErrTracker.TracesSampleRate < 0 || c.ErrTracker.TracesSampleRate > 1.0 {
		return fmt.Errorf("invalid environment variable: %s should be between 0 and 1", "SENTRY_TRACES_SAMPLE_RATE")
	}

	// Mailer
	if c.Mailer.MaxRetries < 0 {
		return fmt.Errorf("invalid environment variable: %s", "MAILER_MAX_RETRIES")
	}

	if c.Mailer.RetryDelayInSeconds < 0 {
		return fmt.Errorf("invalid environment variable: %s", "MAILER_RETRIES_DELAY_IN_SECONDS")
	}

	return nil
}
