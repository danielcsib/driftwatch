package drift_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func Example_loadAndDetect() {
	// Write a temporary desired config file.
	dir, _ := os.MkdirTemp("", "driftwatch-example")
	defer os.RemoveAll(dir)

	content := `{"name":"api","image":"api:v2","replicas":1,"environment":{"PORT":"9090"}}`
	_ = os.WriteFile(filepath.Join(dir, "api.json"), []byte(content), 0644)

	// Load desired config from file.
	desired, err := drift.LoadConfig(filepath.Join(dir, "api.json"))
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Simulate a live config that differs in image.
	live := &drift.ServiceConfig{
		Name:        desired.Name,
		Image:       "api:v1", // stale image
		Replicas:    desired.Replicas,
		Environment: desired.Environment,
	}

	detector := drift.NewDetector()
	result := detector.Detect(desired, live)

	if result.Drifted {
		fmt.Println("drift detected")
	} else {
		fmt.Println("no drift")
	}
	// Output: drift detected
}
