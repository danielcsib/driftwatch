package drift

import (
	"errors"
	"fmt"
	"time"
)

// ScheduleConfig holds raw configuration for a Scheduler, suitable for
// unmarshalling from JSON or YAML.
type ScheduleConfig struct {
	// IntervalSeconds is the polling interval in seconds. Minimum 5.
	IntervalSeconds int `json:"interval_seconds" yaml:"interval_seconds"`
}

// Validate checks that ScheduleConfig values are within acceptable ranges.
func (c ScheduleConfig) Validate() error {
	if c.IntervalSeconds < 5 {
		return errors.New("interval_seconds must be at least 5")
	}
	return nil
}

// ToSchedule converts a validated ScheduleConfig to a Schedule.
func (c ScheduleConfig) ToSchedule() (Schedule, error) {
	if err := c.Validate(); err != nil {
		return Schedule{}, fmt.Errorf("invalid schedule config: %w", err)
	}
	return Schedule{
		Interval: time.Duration(c.IntervalSeconds) * time.Second,
	}, nil
}

// DefaultScheduleConfig returns a ScheduleConfig with sensible defaults.
func DefaultScheduleConfig() ScheduleConfig {
	return ScheduleConfig{
		IntervalSeconds: 30,
	}
}
