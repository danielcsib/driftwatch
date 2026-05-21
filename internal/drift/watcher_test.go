package drift_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/drift"
)

func writeCfgFile(t *testing.T, dir string, cfg drift.ServiceConfig) {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	path := filepath.Join(dir, cfg.Name+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}

func TestWatcher_NoDrift(t *testing.T) {
	dir := t.TempDir()
	cfg := drift.ServiceConfig{
		Name:  "api",
		Image: "api:v1",
		Env:   map[string]string{"PORT": "8080"},
	}
	writeCfgFile(t, dir, cfg)

	detector := drift.NewDetector()
	reporter := drift.NewReporter("text")
	watcher := drift.NewWatcher(
		drift.WatchConfig{
			Interval:  50 * time.Millisecond,
			ConfigDir: dir,
			Desired:   []drift.ServiceConfig{cfg},
		},
		detector,
		reporter,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	go watcher.Run(ctx) //nolint:errcheck

	result := <-watcher.Results()
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if len(result.Reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(result.Reports))
	}
	if result.Reports[0].Drifted {
		t.Errorf("expected no drift, got drift: %v", result.Reports[0].Diffs)
	}
}

func TestWatcher_DetectsDrift(t *testing.T) {
	dir := t.TempDir()
	actual := drift.ServiceConfig{
		Name:  "worker",
		Image: "worker:v2", // deployed version differs
		Env:   map[string]string{"QUEUE": "jobs"},
	}
	writeCfgFile(t, dir, actual)

	desired := drift.ServiceConfig{
		Name:  "worker",
		Image: "worker:v1",
		Env:   map[string]string{"QUEUE": "jobs"},
	}

	detector := drift.NewDetector()
	reporter := drift.NewReporter("text")
	watcher := drift.NewWatcher(
		drift.WatchConfig{
			Interval:  50 * time.Millisecond,
			ConfigDir: dir,
			Desired:   []drift.ServiceConfig{desired},
		},
		detector,
		reporter,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	go watcher.Run(ctx) //nolint:errcheck

	result := <-watcher.Results()
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if !result.Reports[0].Drifted {
		t.Error("expected drift to be detected")
	}
}
