package drift_test

import (
	"os"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

func makeSnapshot(t time.Time, drifted bool) drift.Snapshot {
	results := []drift.DetectResult{
		{
			Service: "api",
			Drifted: drifted,
			Fields:  []string{"image"},
		},
	}
	return drift.Snapshot{Timestamp: t, Results: results}
}

func TestSnapshotStore_SaveAndLatest(t *testing.T) {
	dir := t.TempDir()
	store, err := drift.NewSnapshotStore(dir)
	if err != nil {
		t.Fatalf("NewSnapshotStore: %v", err)
	}

	snap := makeSnapshot(time.Now().UTC(), true)
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Latest()
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if len(got.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(got.Results))
	}
	if got.Results[0].Service != "api" {
		t.Errorf("expected service 'api', got %q", got.Results[0].Service)
	}
}

func TestSnapshotStore_Latest_Empty(t *testing.T) {
	dir := t.TempDir()
	store, _ := drift.NewSnapshotStore(dir)

	_, err := store.Latest()
	if err == nil {
		t.Error("expected error for empty store, got nil")
	}
}

func TestSnapshotStore_List(t *testing.T) {
	dir := t.TempDir()
	store, _ := drift.NewSnapshotStore(dir)

	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 3; i++ {
		snap := makeSnapshot(base.Add(time.Duration(i)*time.Second), i%2 == 0)
		if err := store.Save(snap); err != nil {
			t.Fatalf("Save[%d]: %v", i, err)
		}
	}

	snaps, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(snaps) != 3 {
		t.Errorf("expected 3 snapshots, got %d", len(snaps))
	}
}

func TestSnapshotStore_CreatesDirIfMissing(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/snaps"

	_, err := drift.NewSnapshotStore(dir)
	if err != nil {
		t.Fatalf("expected dir creation, got: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("directory was not created")
	}
}
