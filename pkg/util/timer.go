package util

import "time"

// Timer wraps time.Timer with an explicit end time tracking.
type Timer struct {
	timer *time.Timer
	end   time.Time
}

// NewTimer creates a new Timer instance.
func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		timer: time.NewTimer(duration),
		end:   time.Now().Add(duration),
	}
}

func (s *Timer) C() <-chan time.Time {
	return s.timer.C
}

// Reset restarts the timer with a new duration.
func (s *Timer) Reset(duration time.Duration) {
	s.timer.Reset(duration)
	s.end = time.Now().Add(duration)
}

// Stop stops the timer.
func (s *Timer) Stop() {
	if s.timer != nil {
		s.timer.Stop()
	}
}

// TimeRemaining returns the remaining duration.
func (s *Timer) TimeRemaining() time.Duration {
	remaining := time.Until(s.end)
	if remaining < 0 {
		return 0
	}
	return remaining
}
