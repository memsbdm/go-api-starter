package mocks

import (
	"go-starter/config"
	"go-starter/internal/domain/ports"
	"net/http"
	"time"
)

// ErrorTrackerMock implements the ports.ErrorTracker interface.
// It is not implemented and used in local development and tests
type ErrorTrackerMock struct{}

// NewErrorTrackerMock creates and returns a new ErrorTrackerMock instance.
func NewErrorTrackerMock(_ *config.ErrTracker) *ErrorTrackerMock {
	return &ErrorTrackerMock{}
}

// Handle wraps the provided http.Handler with a middleware for automatic
// error tracking and request monitoring.
func (set *ErrorTrackerMock) Handle(_ http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

// SetUser associates the current scope with user information identified by
// the provided ID and IP address.
func (set *ErrorTrackerMock) SetUser(_, _ string) {}

// CaptureException sends an error and returns the event ID as a string.
func (set *ErrorTrackerMock) CaptureException(_ error) string {
	return ""
}

// AddBreadcrumb adds a new breadcrumb to the current scope with the specified
// message and options. Breadcrumbs track the series of events leading up to an error.
func (set *ErrorTrackerMock) AddBreadcrumb(_ string, _ ports.BreadCrumbOptions) {}

// SetRequest attaches the provided HTTP request to the current scope for
// additional context in error reports.
func (set *ErrorTrackerMock) SetRequest(_ *http.Request) {}

// SetBody attaches the provided request body to the current scope for
// additional context in error reports.
func (set *ErrorTrackerMock) SetBody(_ []byte) {}

// Flush waits for queued events to be sent for the specified duration.
// It should be called before program termination to ensure all events are sent.
func (set *ErrorTrackerMock) Flush(_ time.Duration) {}
