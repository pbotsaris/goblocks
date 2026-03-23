package socketmode

import (
	"math"
	"math/rand/v2"
	"time"
)

// Backoff implements exponential backoff with jitter.
type Backoff struct {
	baseDelay  time.Duration
	maxDelay   time.Duration
	maxJitter  time.Duration
	attempts   int
	stableTime time.Duration
	lastStable time.Time
}

// BackoffConfig configures the backoff behavior.
type BackoffConfig struct {
	BaseDelay  time.Duration // Initial delay (default: 1s)
	MaxDelay   time.Duration // Maximum delay (default: 30s)
	MaxJitter  time.Duration // Maximum random jitter (default: 1s)
	StableTime time.Duration // Time before resetting attempts (default: 60s)
}

// DefaultBackoffConfig returns the default backoff configuration.
func DefaultBackoffConfig() BackoffConfig {
	return BackoffConfig{
		BaseDelay:  1 * time.Second,
		MaxDelay:   30 * time.Second,
		MaxJitter:  1 * time.Second,
		StableTime: 60 * time.Second,
	}
}

// NewBackoff creates a new Backoff with the given configuration.
func NewBackoff(cfg BackoffConfig) *Backoff {
	if cfg.BaseDelay <= 0 {
		cfg.BaseDelay = 1 * time.Second
	}
	if cfg.MaxDelay <= 0 {
		cfg.MaxDelay = 30 * time.Second
	}
	if cfg.MaxJitter < 0 {
		cfg.MaxJitter = 1 * time.Second
	}
	if cfg.StableTime <= 0 {
		cfg.StableTime = 60 * time.Second
	}
	return &Backoff{
		baseDelay:  cfg.BaseDelay,
		maxDelay:   cfg.MaxDelay,
		maxJitter:  cfg.MaxJitter,
		stableTime: cfg.StableTime,
	}
}

// NextDelay returns the next backoff delay and increments the attempt counter.
// Formula: min(baseDelay * 2^attempt + jitter, maxDelay)
func (b *Backoff) NextDelay() time.Duration {
	// Calculate exponential delay
	delay := float64(b.baseDelay) * math.Pow(2, float64(b.attempts))

	// Cap at max delay
	if delay > float64(b.maxDelay) {
		delay = float64(b.maxDelay)
	}

	// Add jitter (0 to maxJitter)
	if b.maxJitter > 0 {
		jitter := time.Duration(rand.Int64N(int64(b.maxJitter)))
		delay += float64(jitter)
	}

	// Increment attempts
	b.attempts++

	return time.Duration(delay)
}

// Attempts returns the current number of consecutive failed attempts.
func (b *Backoff) Attempts() int {
	return b.attempts
}

// Reset resets the attempt counter to zero.
func (b *Backoff) Reset() {
	b.attempts = 0
	b.lastStable = time.Time{}
}

// MarkConnected should be called when a connection is successfully established.
// It records the time for stability tracking.
func (b *Backoff) MarkConnected() {
	b.lastStable = time.Now()
}

// CheckStable checks if the connection has been stable long enough to reset.
// Returns true if the backoff counter was reset.
func (b *Backoff) CheckStable() bool {
	if b.lastStable.IsZero() {
		return false
	}
	if time.Since(b.lastStable) >= b.stableTime {
		b.Reset()
		return true
	}
	return false
}
