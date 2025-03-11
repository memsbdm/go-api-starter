package errtracker

import (
	"go-starter/internal/domain/ports"
	"net/http"
	"time"
)

// ErrTrackerAdapterMock implements the ports.ErrTrackerAdapter interface.
// It is not implemented and used in local development and tests
type ErrTrackerAdapterMock struct{}

// NewErrTrackerAdapterMock creates and returns a new ErrTrackerAdapterMock instance.
func NewErrTrackerAdapterMock() *ErrTrackerAdapterMock {
	return &ErrTrackerAdapterMock{}
}

// Handle wraps the provided http.Handler with a middleware for automatic
// error tracking and request monitoring.
func (mock *ErrTrackerAdapterMock) Handle(handler http.Handler) http.Handler {
	return handler
}

// SetUser associates the current scope with user information identified by
// the provided ID and IP address.
func (mock *ErrTrackerAdapterMock) SetUser(_, _ string) {}

// CaptureException sends an error and returns the event ID as a string.
func (mock *ErrTrackerAdapterMock) CaptureException(_ error) string {
	return ""
}

// AddBreadcrumb adds a new breadcrumb to the current scope with the specified
// message and options. Breadcrumbs track the series of events leading up to an error.
func (mock *ErrTrackerAdapterMock) AddBreadcrumb(_ string, _ ports.BreadCrumbOptions) {}

// SetRequest attaches the provided HTTP request to the current scope for
// additional context in error reports.
func (mock *ErrTrackerAdapterMock) SetRequest(_ *http.Request) {}

// SetBody attaches the provided request body to the current scope for
// additional context in error reports.
func (mock *ErrTrackerAdapterMock) SetBody(_ []byte) {}

// Flush waits for queued events to be sent for the specified duration.
// It should be called before program termination to ensure all events are sent.
func (mock *ErrTrackerAdapterMock) Flush(_ time.Duration) {}
