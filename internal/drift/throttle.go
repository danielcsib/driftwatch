package drift

import (
	"sync"
	"time"
)

// ThrottleConfig controls how often alerts fire for the same service.
type ThrottleConfig struct {
	// Cooldown is the minimum duration between alerts for the same service.
	Cooldown time.Duration
}

// DefaultThrottleConfig returns a ThrottleConfig with sensible defaults.
func DefaultThrottleConfig() ThrottleConfig {
	return ThrottleConfig{
		Cooldown: 5 * time.Minute,
	}
}

// Throttle suppresses repeated alerts for the same service within a cooldown
// window.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastSent map[string]time.Time
	now      func() time.Time
}

// NewThrottle creates a Throttle using the provided config.
func NewThrottle(cfg ThrottleConfig) *Throttle {
	return &Throttle{
		cooldown: cfg.Cooldown,
		lastSent: make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if an alert for the given service should be sent,
// and records the send time if so.
func (t *Throttle) Allow(service string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.lastSent[service]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}
	t.lastSent[service] = now
	return true
}

// Reset clears the recorded send time for a service, allowing the next
// alert to pass through immediately.
func (t *Throttle) Reset(service string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSent, service)
}

// LastSent returns the time the last alert was sent for a service and
// whether a record exists.
func (t *Throttle) LastSent(service string) (time.Time, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t2, ok := t.lastSent[service]
	return t2, ok
}
