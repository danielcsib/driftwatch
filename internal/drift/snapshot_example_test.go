package drift_test

import (
	"fmt"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

func Example_snapshotSaveAndList() {
	dir, _ := os.MkdirTemp("", "driftwatch-snap-*")
	defer os.RemoveAll(dir)

	store, _ := drift.NewSnapshotStore(dir)

	snap := drift.Snapshot{
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Results: []drift.DetectResult{
			{Service: "worker", Drifted: true, Fields: []string{"image"}},
			{Service: "api", Drifted: false},
		},
	}

	if err := store.Save(snap); err != nil {
		fmt.Println("save error:", err)
		return
	}

	snaps, err := store.List()
	if err != nil {
		fmt.Println("list error:", err)
		return
	}

	for _, s := range snaps {
		fmt.Printf("snapshot at %s: %d results\n",
			s.Timestamp.Format(time.RFC3339), len(s.Results))
	}
	// Output:
	// snapshot at 2024-06-01T12:00:00Z: 2 results
}
