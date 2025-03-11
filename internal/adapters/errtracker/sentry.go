package errtracker

import (
	"fmt"
	"go-starter/config"
	"go-starter/internal/domain/ports"
	"log/slog"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

// SentryAdapter implements ports.ErrTrackerAdapter interface and provides integration
// with Sentry error monitoring service.
type SentryAdapter struct {
	cfg *config.Container
}

// New creates a new instance of SentryAdapter.
func New(cfg *config.Container) *SentryAdapter {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.ErrTracker.DSN,
		TracesSampleRate: cfg.ErrTracker.TracesSampleRate,
	}); err != nil {
		slog.Error(fmt.Sprintf("Sentry initialization failed: %v\n", err))
	}

	return &SentryAdapter{
		cfg: cfg,
	}
}

// Handle wraps the provided http.Handler with Sentry middleware for automatic
// error tracking and request monitoring.
func (sa *SentryAdapter) Handle(handler http.Handler) http.Handler {
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	return sentryHandler.Handle(handler)
}

// SetUser associates the current scope with user information identified by
// the provided ID and IP address.
func (sa *SentryAdapter) SetUser(id, ipAddr string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID:        id,
			IPAddress: ipAddr,
		})
	})
}

// CaptureException sends an error to Sentry and returns the event ID as a string.
func (sa *SentryAdapter) CaptureException(err error) string {
	event := &sentry.Event{
		Environment: sa.cfg.Application.Env,
		Exception:   []sentry.Exception{{Value: err.Error()}},
	}

	eventID := sentry.CaptureEvent(event)

	return string(*eventID)
}

// AddBreadcrumb adds a new breadcrumb to the current scope with the specified
// message and options. Breadcrumbs track the series of events leading up to an error.
func (sa *SentryAdapter) AddBreadcrumb(message string, options ports.BreadCrumbOptions) {
	level := sentry.LevelError
	if options.Level != "" {
		level = mapDomainSentryLevel(options.Level)
	}

	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Timestamp: time.Now(),
		Message:   message,
		Level:     level,
		Category:  options.Category,
		Data:      options.Data,
	})
}

// SetRequest attaches the provided HTTP request to the current scope for
// additional context in error reports.
func (sa *SentryAdapter) SetRequest(r *http.Request) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetRequest(r)
	})
}

// SetBody attaches the provided request body to the current scope for
// additional context in error reports.
func (sa *SentryAdapter) SetBody(body []byte) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetRequestBody(body)
	})
}

// Flush waits for queued events to be sent to Sentry for the specified duration.
// It should be called before program termination to ensure all events are sent.
func (sa *SentryAdapter) Flush(duration time.Duration) {
	sentry.Flush(duration)
}

// mapDomainSentryLevel converts internal error tracking levels to corresponding
// Sentry levels. It defaults to LevelError if the level is not recognized.
func mapDomainSentryLevel(level ports.ErrTrackerLevel) sentry.Level {
	switch level {
	case ports.LevelError:
		return sentry.LevelError
	case ports.LevelWarning:
		return sentry.LevelWarning
	case ports.LevelInfo:
		return sentry.LevelInfo
	case ports.LevelFatal:
		return sentry.LevelFatal
	}
	return sentry.LevelError
}
