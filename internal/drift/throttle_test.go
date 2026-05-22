package drift

import (
	"testing"
	"time"
)

func TestThrottle_Allow_FirstCall(t *testing.T) {
	th := NewThrottle(DefaultThrottleConfig())
	if !th.Allow("svc-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestThrottle_Allow_WithinCooldown(t *testing.T) {
	now := time.Now()
	th := NewThrottle(ThrottleConfig{Cooldown: 10 * time.Minute})
	th.now = func() time.Time { return now }

	if !th.Allow("svc-a") {
		t.Fatal("expected first call to be allowed")
	}
	// Advance time by less than the cooldown.
	th.now = func() time.Time { return now.Add(3 * time.Minute) }
	if th.Allow("svc-a") {
		t.Fatal("expected second call within cooldown to be suppressed")
	}
}

func TestThrottle_Allow_AfterCooldown(t *testing.T) {
	now := time.Now()
	th := NewThrottle(ThrottleConfig{Cooldown: 5 * time.Minute})
	th.now = func() time.Time { return now }

	th.Allow("svc-b")
	th.now = func() time.Time { return now.Add(6 * time.Minute) }
	if !th.Allow("svc-b") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestThrottle_Allow_IndependentServices(t *testing.T) {
	th := NewThrottle(DefaultThrottleConfig())
	if !th.Allow("svc-x") {
		t.Fatal("expected svc-x to be allowed")
	}
	if !th.Allow("svc-y") {
		t.Fatal("expected svc-y to be allowed independently")
	}
}

func TestThrottle_Reset(t *testing.T) {
	now := time.Now()
	th := NewThrottle(ThrottleConfig{Cooldown: 10 * time.Minute})
	th.now = func() time.Time { return now }

	th.Allow("svc-c")
	th.Reset("svc-c")
	// After reset the next call should be allowed even within cooldown.
	if !th.Allow("svc-c") {
		t.Fatal("expected allow after reset")
	}
}

func TestThrottle_LastSent(t *testing.T) {
	now := time.Now()
	th := NewThrottle(DefaultThrottleConfig())
	th.now = func() time.Time { return now }

	if _, ok := th.LastSent("svc-d"); ok {
		t.Fatal("expected no last-sent record before first allow")
	}
	th.Allow("svc-d")
	got, ok := th.LastSent("svc-d")
	if !ok {
		t.Fatal("expected last-sent record after allow")
	}
	if !got.Equal(now) {
		t.Fatalf("expected last-sent %v, got %v", now, got)
	}
}

func TestDefaultThrottleConfig(t *testing.T) {
	cfg := DefaultThrottleConfig()
	if cfg.Cooldown != 5*time.Minute {
		t.Fatalf("expected default cooldown 5m, got %v", cfg.Cooldown)
	}
}
