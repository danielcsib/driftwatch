package drift_test

import (
	"fmt"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Example_historyRecordAndQuery shows how to record drift events and
// query the latest result.
func Example_historyRecordAndQuery() {
	h := drift.NewHistory(50)

	// Simulate a clean check.
	h.Record([]drift.DriftResult{
		{Service: "api", Drifted: false},
	})

	// Simulate a drifted check.
	h.Record([]drift.DriftResult{
		{Service: "api", Drifted: true, Details: "image tag changed"},
	})

	ev, ok := h.Latest()
	if !ok {
		fmt.Println("no events")
		return
	}

	if ev.HasDrift {
		fmt.Println("drift detected in latest check")
	} else {
		fmt.Println("no drift in latest check")
	}

	fmt.Printf("total events recorded: %d\n", h.Len())

	// Output:
	// drift detected in latest check
	// total events recorded: 2
}
