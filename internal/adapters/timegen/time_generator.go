package timegen

import "time"

// RealTimeGenerator implements ports.TimeGenerator interface for real time management.
type RealTimeGenerator struct{}

// NewRealTimeGenerator creates a new instance of RealTimeGenerator.
func NewRealTimeGenerator() *RealTimeGenerator {
	return &RealTimeGenerator{}
}

// Now returns the current real time.
func (rtg *RealTimeGenerator) Now() time.Time {
	return time.Now()
}

// Advance is not implemented for RealTimeGenerator.
// This method is meant to modify time for testing purposes and does nothing here.
func (rtg *RealTimeGenerator) Advance(_ time.Duration) {}

// FakeTimeGenerator implements ports.TimeGenerator interface for testing with fake time.
type FakeTimeGenerator struct {
	time time.Time
}

// NewFakeTimeGenerator creates a new instance of FakeTimeGenerator with a specified initial time.
func NewFakeTimeGenerator(time time.Time) *FakeTimeGenerator {
	return &FakeTimeGenerator{
		time: time,
	}
}

// Now returns the specified fake time.
func (ftg *FakeTimeGenerator) Now() time.Time {
	return ftg.time
}

// Advance adds a duration to the current fake time.
func (ftg *FakeTimeGenerator) Advance(d time.Duration) {
	ftg.time = ftg.time.Add(d)
}
