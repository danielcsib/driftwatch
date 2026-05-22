package drift

import (
	"context"
	"log"
	"time"
)

// Schedule defines how often drift checks are run.
type Schedule struct {
	Interval time.Duration
}

// Scheduler runs drift checks on a fixed interval.
type Scheduler struct {
	schedule Schedule
	watcher  *Watcher
	alerter  Alerter
	logger   *log.Logger
}

// NewScheduler creates a Scheduler that checks for drift at the given interval.
func NewScheduler(s Schedule, w *Watcher, a Alerter, logger *log.Logger) *Scheduler {
	return &Scheduler{
		schedule: s,
		watcher:  w,
		alerter:  a,
		logger:   logger,
	}
}

// Run starts the scheduler loop, blocking until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.schedule.Interval)
	defer ticker.Stop()

	s.logger.Printf("scheduler: starting drift checks every %s", s.schedule.Interval)

	for {
		select {
		case <-ctx.Done():
			s.logger.Println("scheduler: shutting down")
			return
		case <-ticker.C:
			s.runOnce(ctx)
		}
	}
}

func (s *Scheduler) runOnce(ctx context.Context) {
	results, err := s.watcher.Check()
	if err != nil {
		s.logger.Printf("scheduler: check error: %v", err)
		return
	}

	for _, r := range results {
		if r.Drifted {
			alert := BuildAlert(r)
			if err := s.alerter.Send(ctx, alert); err != nil {
				s.logger.Printf("scheduler: alert error for %s: %v", r.ServiceName, err)
			}
		}
	}
}
