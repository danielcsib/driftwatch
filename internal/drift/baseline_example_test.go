package drift_test

import (
	"fmt"
	"os"

	"github.com/example/driftwatch/internal/drift"
)

func Example_baselineCompare() {
	dir, _ := os.MkdirTemp("", "baseline-example")
	defer os.RemoveAll(dir)

	store, _ := drift.NewBaselineStore(dir)

	original := []drift.ServiceConfig{
		{Name: "api", Image: "api:1.0", Env: map[string]string{"PORT": "8080"}},
	}
	_ = store.Save("prod", original)

	live := []drift.ServiceConfig{
		{Name: "api", Image: "api:2.0", Env: map[string]string{"PORT": "8080"}},
	}

	b, _ := store.Load("prod")
	diffs := drift.CompareToBaseline(b, live)

	for _, d := range diffs {
		fmt.Printf("service=%s field=%s baseline=%s live=%s\n",
			d.Service, d.Field, d.Baseline, d.Live)
	}
	// Output:
	// service=api field=image baseline=api:1.0 live=api:2.0
}
