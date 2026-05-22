package drift_test

import (
	"fmt"

	"github.com/yourorg/driftwatch/internal/drift"
)

func Example_policyEvaluate() {
	p := drift.NewPolicy([]drift.PolicyRule{
		{Service: "api", Field: "image", Action: drift.PolicyActionIgnore},
		{Service: "*", Field: "replicas", Action: drift.PolicyActionFail},
	})

	fmt.Println(p.Evaluate("api", "image"))    // ignore
	fmt.Println(p.Evaluate("worker", "replicas")) // fail
	fmt.Println(p.Evaluate("db", "env"))          // alert (default)

	// Output:
	// ignore
	// fail
	// alert
}

func Example_applyPolicy() {
	p := drift.NewPolicy([]drift.PolicyRule{
		{Service: "api", Field: "image", Action: drift.PolicyActionIgnore},
	})

	results := []drift.DriftResult{
		{Service: "api", Field: "image", Drifted: true, Desired: "v1", Actual: "v2"},
		{Service: "api", Field: "env", Drifted: true, Desired: "prod", Actual: "staging"},
	}

	filtered := drift.ApplyPolicy(results, p)
	fmt.Println(len(filtered)) // image rule ignored, only env remains
	fmt.Println(filtered[0].Field)

	// Output:
	// 1
	// env
}
