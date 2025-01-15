package ports

import "time"

// TimeGenerator is an interface for interacting with time.
type TimeGenerator interface {
	// Now returns the current time.
	// In a real implementation, this typically calls time.Now().
	// In test implementations, this can return a fixed time or a time that changes based on test scenarios.
	Now() time.Time

	// Advance changes the current time by adding the specified duration.
	Advance(d time.Duration)
}
