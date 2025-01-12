package ports

import "time"

// TimeGenerator is an interface for interacting with time
type TimeGenerator interface {
	// Now returns the current time
	Now() time.Time
	// Advance changes the current time during tests
	Advance(d time.Duration)
}
