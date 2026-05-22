package drift_test

import (
	"fmt"

	"github.com/yourusername/driftwatch/internal/drift"
)

func Example_summaryFromResults() {
	results := []drift.DetectResult{
		{Service: "api", Drifted: true, DriftedFields: []string{"image"}},
		{Service: "worker", Drifted: false},
	}

	s := drift.NewSummary(results)

	fmt.Printf("Total: %d, Drifted: %d, Clean: %d\n",
		s.TotalChecked, s.DriftedCount, s.CleanCount)

	for _, svc := range s.Services {
		if svc.Drifted {
			fmt.Printf("Service %q drifted on: %v\n", svc.Name, svc.Fields)
		} else {
			fmt.Printf("Service %q is clean\n", svc.Name)
		}
	}

	// Output:
	// Total: 2, Drifted: 1, Clean: 1
	// Service "api" drifted on: [image]
	// Service "worker" is clean
}
