package logger

import (
	"go-starter/config"
	"log/slog"
	"os"
)

// New defines the logger specifications based on the application environment.
// It initializes a new logger and sets it as the default logger for the application.
// In production, it uses a JSON format for logging; otherwise, it uses a plain text format.
func New(appCfg *config.App) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if appCfg.Env != config.EnvDevelopment {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	slog.SetDefault(logger)
}
