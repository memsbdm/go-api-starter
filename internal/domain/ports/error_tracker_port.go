package ports

import (
	"net/http"
	"time"
)

// ErrorTracker is an interface for interacting with error monitoring business logic.
type ErrorTracker interface {
	// SetUser associates the current scope with user information identified by
	// the provided ID and IP address.
	SetUser(id, ipAddr string)
	// SetRequest attaches the provided HTTP request to the current scope for
	// additional context in error reports.
	SetRequest(r *http.Request)
	// SetBody attaches the provided request body to the current scope for
	// additional context in error reports.
	SetBody(body []byte)
	// Handle wraps the provided http.Handler with a middleware for automatic
	// error tracking and request monitoring.
	Handle(handler http.Handler) http.Handler
	// CaptureException sends an error and returns the event ID as a string.
	CaptureException(err error) string
	// AddBreadcrumb adds a new breadcrumb to the current scope with the specified
	// message and options. Breadcrumbs track the series of events leading up to an error.
	AddBreadcrumb(message string, options BreadCrumbOptions)
	// Flush waits for queued events to be sent for the specified duration.
	// It should be called before program termination to ensure all events are sent.
	Flush(duration time.Duration)
}

// ErrorTrackerLevel represents the severity level of an error or event
// in the error tracking system.
type ErrorTrackerLevel string

const (
	LevelFatal   ErrorTrackerLevel = "fatal"
	LevelError   ErrorTrackerLevel = "error"
	LevelWarning ErrorTrackerLevel = "warning"
	LevelInfo    ErrorTrackerLevel = "info"
)

// BreadCrumbOptions contains configuration options for creating a breadcrumb
// in the error tracking system.
type BreadCrumbOptions struct {
	Type     string
	Category string
	Level    ErrorTrackerLevel
	Data     map[string]interface{}
}
