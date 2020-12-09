package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Timer implements an
type Timer struct {
	// Updated is used to notify UI about changes in the timer.
	Updated chan struct{}

	// mu locks the state such that it can be modified and accessed
	// from multiple goroutines.
	mu       sync.Mutex
	start    time.Time     // start corresponds to when the timer was started.
	now      time.Time     // now corresponds to the last updated time.
	duration time.Duration // duration is the maximum progress.
}

// NewTimer creates a new timer with the specified timer.
func NewTimer(initialDuration time.Duration) *Timer {
	return &Timer{
		Updated:  make(chan struct{}),
		duration: initialDuration,
	}
}

// Start the timer goroutine and return a cancel func that
// that can be used to stop it.
func (t *Timer) Start() context.CancelFunc {
	// initialize the timer state.
	now := time.Now()
	t.now = now
	t.start = now

	// we use done to signal stopping the goroutine.
	// a context.Context could be also used.
	done := make(chan struct{})
	go t.run(done)
	return func() { close(done) }
}

// run is the main loop for the timer.
func (t *Timer) run(done chan struct{}) {
	// we use a time.Ticker to update the state,
	// in many cases, this could be a network access instead.
	tick := time.NewTicker(50 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case now := <-tick.C:
			t.update(now)
		case <-done:
			return
		}
	}
}

// invalidate sends a signal to the UI that
// the internal state has changed.
func (t *Timer) invalidate() {
	// we use a non-blocking send, that way the Timer
	// can continue updating internally.
	select {
	case t.Updated <- struct{}{}:
	default:
	}
}

func (t *Timer) update(now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	previousNow := t.now
	t.now = now

	// first check whether we have not exceeded the duration.
	// in that case the progress advanced and we need to notify
	// about a change.
	progressAfter := t.now.Sub(t.start)
	if progressAfter <= t.duration {
		t.invalidate()
		return
	}

	// when we had progressed beyond the duration we also
	// need to update the first time it happens.
	progressBefore := previousNow.Sub(t.start)
	if progressBefore <= t.duration {
		t.invalidate()
		return
	}
}

// Reset resets timer to the last know time.
func (t *Timer) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.start = t.now
	t.invalidate()
}

// SetDuration changes the duration of the timer.
func (t *Timer) SetDuration(duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.duration == duration {
		return
	}
	t.duration = duration
	t.invalidate()
}

// Info returns the latest know info about the timer.
func (t *Timer) Info() (info Info) {
	t.mu.Lock()
	defer t.mu.Unlock()

	info.Progress = t.now.Sub(t.start)
	info.Duration = t.duration
	if info.Progress > info.Duration {
		info.Progress = info.Duration
	}
	return info
}

// Info is the information about the timer.
type Info struct {
	Progress time.Duration
	Duration time.Duration
}

// ProgressString returns the progress formatted as seconds.
func (info *Info) ProgressString() string {
	return fmt.Sprintf("%.1fs", info.Progress.Seconds())
}
