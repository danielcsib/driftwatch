package drift_test

import (
	"fmt"

	"github.com/yourusername/driftwatch/internal/drift"
)

// Example_summaryFromResults demonstrates how to build a drift summary
// from a slice of DetectResult values and iterate over the service details.
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

// Example_summaryNoDrift demonstrates that a summary with no drifted services
// reports a DriftedCount of zero.
func Example_summaryNoDrift() {
	results := []drift.DetectResult{
		{Service: "api", Drifted: false},
		{Service: "worker", Drifted: false},
	}

	s := drift.NewSummary(results)

	fmt.Printf("Total: %d, Drifted: %d, Clean: %d\n",
		s.TotalChecked, s.DriftedCount, s.CleanCount)

	// Output:
	// Total: 2, Drifted: 0, Clean: 2
}
