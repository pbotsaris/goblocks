package socketmode

import (
	"testing"
	"time"
)

func TestBackoff_NextDelay(t *testing.T) {
	t.Run("returns base delay on first attempt", func(t *testing.T) {
		b := NewBackoff(BackoffConfig{
			BaseDelay: 1 * time.Second,
			MaxDelay:  30 * time.Second,
			MaxJitter: 0, // No jitter for predictable tests
		})

		delay := b.NextDelay()
		if delay != 1*time.Second {
			t.Errorf("got delay %v, want 1s", delay)
		}
	})

	t.Run("exponential growth", func(t *testing.T) {
		b := NewBackoff(BackoffConfig{
			BaseDelay: 1 * time.Second,
			MaxDelay:  30 * time.Second,
			MaxJitter: 0,
		})

		expected := []time.Duration{
			1 * time.Second,  // 1 * 2^0
			2 * time.Second,  // 1 * 2^1
			4 * time.Second,  // 1 * 2^2
			8 * time.Second,  // 1 * 2^3
			16 * time.Second, // 1 * 2^4
			30 * time.Second, // capped at max
			30 * time.Second, // still capped
		}

		for i, want := range expected {
			got := b.NextDelay()
			if got != want {
				t.Errorf("attempt %d: got %v, want %v", i, got, want)
			}
		}
	})

	t.Run("never exceeds maxDelay", func(t *testing.T) {
		b := NewBackoff(BackoffConfig{
			BaseDelay: 10 * time.Second,
			MaxDelay:  15 * time.Second,
			MaxJitter: 0,
		})

		// First delay: 10s
		// Second delay: 20s -> capped to 15s
		b.NextDelay()
		delay := b.NextDelay()
		if delay > 15*time.Second {
			t.Errorf("got delay %v, want <= 15s", delay)
		}
	})

	t.Run("jitter is within bounds", func(t *testing.T) {
		b := NewBackoff(BackoffConfig{
			BaseDelay: 1 * time.Second,
			MaxDelay:  30 * time.Second,
			MaxJitter: 500 * time.Millisecond,
		})

		// Run multiple times to test jitter randomness
		for i := 0; i < 100; i++ {
			b.Reset()
			delay := b.NextDelay()
			if delay < 1*time.Second || delay > 1500*time.Millisecond {
				t.Errorf("delay %v out of expected range [1s, 1.5s]", delay)
			}
		}
	})

	t.Run("increments attempts", func(t *testing.T) {
		b := NewBackoff(BackoffConfig{
			BaseDelay: 1 * time.Second,
			MaxDelay:  30 * time.Second,
			MaxJitter: 0,
		})

		if b.Attempts() != 0 {
			t.Errorf("initial attempts: got %d, want 0", b.Attempts())
		}

		b.NextDelay()
		if b.Attempts() != 1 {
			t.Errorf("after first delay: got %d, want 1", b.Attempts())
		}

		b.NextDelay()
		if b.Attempts() != 2 {
			t.Errorf("after second delay: got %d, want 2", b.Attempts())
		}
	})
}

func TestBackoff_Reset(t *testing.T) {
	t.Run("resets attempts to zero", func(t *testing.T) {
		b := NewBackoff(BackoffConfig{
			BaseDelay: 1 * time.Second,
			MaxDelay:  30 * time.Second,
			MaxJitter: 0,
		})

		b.NextDelay()
		b.NextDelay()
		b.NextDelay()

		if b.Attempts() != 3 {
			t.Errorf("before reset: got %d, want 3", b.Attempts())
		}

		b.Reset()

		if b.Attempts() != 0 {
			t.Errorf("after reset: got %d, want 0", b.Attempts())
		}
	})

	t.Run("NextDelay starts from base after reset", func(t *testing.T) {
		b := NewBackoff(BackoffConfig{
			BaseDelay: 1 * time.Second,
			MaxDelay:  30 * time.Second,
			MaxJitter: 0,
		})

		// Build up some attempts
		b.NextDelay() // 1s
		b.NextDelay() // 2s
		b.NextDelay() // 4s

		b.Reset()

		delay := b.NextDelay()
		if delay != 1*time.Second {
			t.Errorf("after reset: got %v, want 1s", delay)
		}
	})
}

func TestBackoff_Stability(t *testing.T) {
	t.Run("MarkConnected records time", func(t *testing.T) {
		b := NewBackoff(BackoffConfig{
			BaseDelay:  1 * time.Second,
			MaxDelay:   30 * time.Second,
			StableTime: 100 * time.Millisecond,
		})

		// Before marking connected, CheckStable should return false
		if b.CheckStable() {
			t.Error("CheckStable should return false before MarkConnected")
		}

		b.MarkConnected()

		// Immediately after, still not stable
		if b.CheckStable() {
			t.Error("CheckStable should return false immediately after MarkConnected")
		}
	})

	t.Run("CheckStable returns true after stableTime", func(t *testing.T) {
		b := NewBackoff(BackoffConfig{
			BaseDelay:  1 * time.Second,
			MaxDelay:   30 * time.Second,
			StableTime: 50 * time.Millisecond,
		})

		b.NextDelay() // Build up attempts
		b.NextDelay()

		if b.Attempts() != 2 {
			t.Fatalf("expected 2 attempts, got %d", b.Attempts())
		}

		b.MarkConnected()

		// Wait for stability
		time.Sleep(60 * time.Millisecond)

		if !b.CheckStable() {
			t.Error("CheckStable should return true after stableTime")
		}

		// Attempts should be reset
		if b.Attempts() != 0 {
			t.Errorf("after stable: got %d attempts, want 0", b.Attempts())
		}
	})

	t.Run("CheckStable returns false before stableTime", func(t *testing.T) {
		b := NewBackoff(BackoffConfig{
			BaseDelay:  1 * time.Second,
			MaxDelay:   30 * time.Second,
			StableTime: 1 * time.Hour, // Very long
		})

		b.MarkConnected()

		if b.CheckStable() {
			t.Error("CheckStable should return false before stableTime")
		}
	})
}

func TestBackoff_Config(t *testing.T) {
	t.Run("uses default values for zero config", func(t *testing.T) {
		b := NewBackoff(BackoffConfig{})

		delay := b.NextDelay()
		// Default base delay is 1s, with jitter up to 1s
		if delay < 1*time.Second || delay > 2*time.Second {
			t.Errorf("delay %v out of default range [1s, 2s]", delay)
		}
	})

	t.Run("respects custom config", func(t *testing.T) {
		b := NewBackoff(BackoffConfig{
			BaseDelay: 100 * time.Millisecond,
			MaxDelay:  500 * time.Millisecond,
			MaxJitter: 0,
		})

		delay := b.NextDelay()
		if delay != 100*time.Millisecond {
			t.Errorf("got %v, want 100ms", delay)
		}

		// 200ms, 400ms, then capped
		b.NextDelay()
		b.NextDelay()
		delay = b.NextDelay()
		if delay != 500*time.Millisecond {
			t.Errorf("got %v, want 500ms (capped)", delay)
		}
	})

	t.Run("DefaultBackoffConfig returns sensible defaults", func(t *testing.T) {
		cfg := DefaultBackoffConfig()

		if cfg.BaseDelay != 1*time.Second {
			t.Errorf("BaseDelay: got %v, want 1s", cfg.BaseDelay)
		}
		if cfg.MaxDelay != 30*time.Second {
			t.Errorf("MaxDelay: got %v, want 30s", cfg.MaxDelay)
		}
		if cfg.MaxJitter != 1*time.Second {
			t.Errorf("MaxJitter: got %v, want 1s", cfg.MaxJitter)
		}
		if cfg.StableTime != 60*time.Second {
			t.Errorf("StableTime: got %v, want 60s", cfg.StableTime)
		}
	})
}
