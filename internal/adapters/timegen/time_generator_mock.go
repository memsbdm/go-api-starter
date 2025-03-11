package timegen

import "time"

// TimeGeneratorMock implements ports.TimeGenerator interface for testing with fake time.
type TimeGeneratorMock struct {
	time time.Time
}

// NewTimeGeneratorMock creates a new instance of TimeGeneratorMock with a specified initial time.
func NewTimeGeneratorMock(time time.Time) *TimeGeneratorMock {
	return &TimeGeneratorMock{
		time: time,
	}
}

// Now returns the specified fake time.
func (mock *TimeGeneratorMock) Now() time.Time {
	return mock.time
}

// Advance adds a duration to the current fake time.
func (mock *TimeGeneratorMock) Advance(d time.Duration) {
	mock.time = mock.time.Add(d)
}
