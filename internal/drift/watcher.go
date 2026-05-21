package drift

import (
	"context"
	"log"
	"time"
)

// WatchConfig holds configuration for the drift watcher.
type WatchConfig struct {
	Interval  time.Duration
	ConfigDir string
	Desired   []ServiceConfig
}

// WatchResult holds the outcome of a single watch cycle.
type WatchResult struct {
	Timestamp time.Time
	Reports   []DriftReport
	Err       error
}

// Watcher periodically checks for config drift and emits results.
type Watcher struct {
	cfg      WatchConfig
	detector *Detector
	reporter *Reporter
	results  chan WatchResult
}

// NewWatcher creates a new Watcher with the given configuration.
func NewWatcher(cfg WatchConfig, detector *Detector, reporter *Reporter) *Watcher {
	return &Watcher{
		cfg:      cfg,
		detector: detector,
		reporter: reporter,
		results:  make(chan WatchResult, 8),
	}
}

// Results returns the channel on which watch cycle outcomes are published.
func (w *Watcher) Results() <-chan WatchResult {
	return w.results
}

// Run starts the watch loop and blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	// Run an immediate check before waiting for the first tick.
	w.runCycle()

	for {
		select {
		case <-ticker.C:
			w.runCycle()
		case <-ctx.Done():
			close(w.results)
			return ctx.Err()
		}
	}
}

func (w *Watcher) runCycle() {
	result := WatchResult{Timestamp: time.Now()}

	actual, err := LoadConfigDir(w.cfg.ConfigDir)
	if err != nil {
		result.Err = err
		w.emit(result)
		log.Printf("[watcher] failed to load configs: %v", err)
		return
	}

	for _, desired := range w.cfg.Desired {
		report := w.detector.Detect(desired, actual)
		result.Reports = append(result.Reports, report)
	}

	w.emit(result)
}

func (w *Watcher) emit(r WatchResult) {
	select {
	case w.results <- r:
	default:
		log.Println("[watcher] result channel full, dropping cycle result")
	}
}
