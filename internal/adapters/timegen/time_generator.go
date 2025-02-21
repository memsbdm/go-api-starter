package timegen

import "time"

// TimeGenerator implements ports.TimeGenerator interface for real time management.
type TimeGenerator struct{}

// NewTimeGenerator creates a new instance of TimeGenerator.
func NewTimeGenerator() *TimeGenerator {
	return &TimeGenerator{}
}

// Now returns the current real time.
func (tg *TimeGenerator) Now() time.Time {
	return time.Now()
}
