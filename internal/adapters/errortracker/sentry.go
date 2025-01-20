package errortracker

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"go-starter/config"
	"go-starter/internal/domain/ports"
	"log/slog"
	"net/http"
	"time"
)

// SentryErrorTracker implements ports.ErrorTracker interface and provides integration
// with Sentry error monitoring service.
type SentryErrorTracker struct{}

// NewSentryErrorTracker creates a new instance of SentryErrorTracker.
func NewSentryErrorTracker(cfg *config.ErrTracker) *SentryErrorTracker {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.DSN,
		TracesSampleRate: cfg.TracesSampleRate,
	}); err != nil {
		slog.Error(fmt.Sprintf("Sentry initialization failed: %v\n", err))
	}

	return &SentryErrorTracker{}
}

// Handle wraps the provided http.Handler with Sentry middleware for automatic
// error tracking and request monitoring.
func (set *SentryErrorTracker) Handle(handler http.Handler) http.Handler {
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	return sentryHandler.Handle(handler)
}

// SetUser associates the current scope with user information identified by
// the provided ID and IP address.
func (set *SentryErrorTracker) SetUser(id, ipAddr string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID:        id,
			IPAddress: ipAddr,
		})
	})
}

// CaptureException sends an error to Sentry and returns the event ID as a string.
func (set *SentryErrorTracker) CaptureException(err error) string {
	eventID := sentry.CaptureException(err)
	return string(*eventID)
}

// AddBreadcrumb adds a new breadcrumb to the current scope with the specified
// message and options. Breadcrumbs track the series of events leading up to an error.
func (set *SentryErrorTracker) AddBreadcrumb(message string, options ports.BreadCrumbOptions) {
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
func (set *SentryErrorTracker) SetRequest(r *http.Request) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetRequest(r)
	})
}

// SetBody attaches the provided request body to the current scope for
// additional context in error reports.
func (set *SentryErrorTracker) SetBody(body []byte) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetRequestBody(body)
	})
}

// Flush waits for queued events to be sent to Sentry for the specified duration.
// It should be called before program termination to ensure all events are sent.
func (set *SentryErrorTracker) Flush(duration time.Duration) {
	sentry.Flush(duration)
}

// mapDomainSentryLevel converts internal error tracking levels to corresponding
// Sentry levels. It defaults to LevelError if the level is not recognized.
func mapDomainSentryLevel(level ports.ErrorTrackerLevel) sentry.Level {
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
