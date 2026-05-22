package drift

import (
	"bytes"
	"context"
	"log"
	"testing"
	"time"
)

func TestScheduler_RunsCheck(t *testing.T) {
	dir := t.TempDir()
	cfg := ServiceConfig{
		Name:  "svc-a",
		Image: "app:1.0",
		Env:   map[string]string{"K": "V"},
	}
	writeCfgFile(t, dir, "svc-a.json", cfg)

	w, err := NewWatcher(dir, cfg)
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}

	var buf bytes.Buffer
	alerter := NewAlerter(&buf)
	logger := log.New(&buf, "", 0)

	sched := NewScheduler(Schedule{Interval: 20 * time.Millisecond}, w, alerter, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		sched.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		// scheduler exited cleanly after ctx cancelled
	case <-time.After(200 * time.Millisecond):
		t.Fatal("scheduler did not stop after context cancellation")
	}
}

func TestScheduler_AlertsOnDrift(t *testing.T) {
	dir := t.TempDir()
	desired := ServiceConfig{
		Name:  "svc-b",
		Image: "app:1.0",
		Env:   map[string]string{"MODE": "prod"},
	}
	// Write a drifted config to disk
	drifted := desired
	drifted.Image = "app:2.0"
	writeCfgFile(t, dir, "svc-b.json", drifted)

	w, err := NewWatcher(dir, desired)
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}

	var buf bytes.Buffer
	alerter := NewAlerter(&buf)
	logger := log.New(&buf, "", 0)

	sched := NewScheduler(Schedule{Interval: 20 * time.Millisecond}, w, alerter, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	sched.Run(ctx)

	if buf.Len() == 0 {
		t.Error("expected alert output for drifted service, got none")
	}
}
