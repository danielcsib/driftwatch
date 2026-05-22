package drift_test

import (
	"os"
	"testing"

	"github.com/example/driftwatch/internal/drift"
)

func makeBaselineConfigs() []drift.ServiceConfig {
	return []drift.ServiceConfig{
		{
			Name:  "api",
			Image: "api:1.0",
			Env:   map[string]string{"PORT": "8080", "LOG": "info"},
		},
		{
			Name:  "worker",
			Image: "worker:2.1",
			Env:   map[string]string{"QUEUE": "default"},
		},
	}
}

func TestBaselineStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, err := drift.NewBaselineStore(dir)
	if err != nil {
		t.Fatalf("NewBaselineStore: %v", err)
	}

	configs := makeBaselineConfigs()
	if err := store.Save("v1", configs); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b, err := store.Load("v1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if b.Label != "v1" {
		t.Errorf("label = %q, want %q", b.Label, "v1")
	}
	if len(b.Configs) != 2 {
		t.Errorf("configs len = %d, want 2", len(b.Configs))
	}
}

func TestBaselineStore_Load_Missing(t *testing.T) {
	dir := t.TempDir()
	store, _ := drift.NewBaselineStore(dir)
	_, err := store.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error loading missing baseline")
	}
}

func TestBaselineStore_List(t *testing.T) {
	dir := t.TempDir()
	store, _ := drift.NewBaselineStore(dir)
	configs := makeBaselineConfigs()
	_ = store.Save("v1", configs)
	_ = store.Save("v2", configs)

	labels, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(labels) != 2 {
		t.Errorf("list len = %d, want 2", len(labels))
	}
}

func TestBaselineStore_CreatesDirIfMissing(t *testing.T) {
	dir := t.TempDir()
	nested := dir + "/a/b/baselines"
	store, err := drift.NewBaselineStore(nested)
	if err != nil {
		t.Fatalf("NewBaselineStore nested: %v", err)
	}
	if _, err := os.Stat(nested); err != nil {
		t.Errorf("dir not created: %v", err)
	}
	_ = store
}
