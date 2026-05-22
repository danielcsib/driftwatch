package drift

import (
	"testing"
	"time"
)

func TestScheduleConfig_Validate_Valid(t *testing.T) {
	c := ScheduleConfig{IntervalSeconds: 10}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestScheduleConfig_Validate_TooShort(t *testing.T) {
	c := ScheduleConfig{IntervalSeconds: 3}
	if err := c.Validate(); err == nil {
		t.Error("expected error for interval < 5, got nil")
	}
}

func TestScheduleConfig_ToSchedule(t *testing.T) {
	c := ScheduleConfig{IntervalSeconds: 60}
	s, err := c.ToSchedule()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Interval != 60*time.Second {
		t.Errorf("expected 60s, got %s", s.Interval)
	}
}

func TestScheduleConfig_ToSchedule_Invalid(t *testing.T) {
	c := ScheduleConfig{IntervalSeconds: 0}
	_, err := c.ToSchedule()
	if err == nil {
		t.Error("expected error for zero interval, got nil")
	}
}

func TestDefaultScheduleConfig(t *testing.T) {
	c := DefaultScheduleConfig()
	if c.IntervalSeconds != 30 {
		t.Errorf("expected default 30s, got %d", c.IntervalSeconds)
	}
	if err := c.Validate(); err != nil {
		t.Errorf("default config should be valid, got %v", err)
	}
}
