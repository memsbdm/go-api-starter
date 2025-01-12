package timegen

import "time"

// RealTimeGenerator implements ports.TimeGenerator
type RealTimeGenerator struct{}

// NewRealTimeGenerator creates a time generator instance
func NewRealTimeGenerator() *RealTimeGenerator {
	return &RealTimeGenerator{}
}

// Now returns the current time
func (rtg *RealTimeGenerator) Now() time.Time {
	return time.Now()
}

// Advance is not implemented, it is adding a duration to the current time only in tests
func (rtg *RealTimeGenerator) Advance(_ time.Duration) {}

// FakeTimeGenerator implements ports.TimeGenerator
type FakeTimeGenerator struct {
	time time.Time
}

// NewFakeTimeGenerator creates a fake time generator instance
func NewFakeTimeGenerator(time time.Time) *FakeTimeGenerator {
	return &FakeTimeGenerator{
		time: time,
	}
}

// Now returns the specified fake time
func (ftg *FakeTimeGenerator) Now() time.Time {
	return ftg.time
}

// Advance is adding a duration to the current time
func (ftg *FakeTimeGenerator) Advance(d time.Duration) {
	ftg.time = ftg.time.Add(d)
}
