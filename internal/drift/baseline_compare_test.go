package drift_test

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
)

func savedBaseline(configs []drift.ServiceConfig) *drift.Baseline {
	return &drift.Baseline{Label: "v1", Configs: configs}
}

func TestCompareToBaseline_NoDrift(t *testing.T) {
	configs := makeBaselineConfigs()
	b := savedBaseline(configs)
	diffs := drift.CompareToBaseline(b, configs)
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d: %+v", len(diffs), diffs)
	}
}

func TestCompareToBaseline_ImageChanged(t *testing.T) {
	base := makeBaselineConfigs()
	live := makeBaselineConfigs()
	live[0].Image = "api:2.0"

	diffs := drift.CompareToBaseline(savedBaseline(base), live)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Field != "image" || diffs[0].Service != "api" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestCompareToBaseline_EnvChanged(t *testing.T) {
	base := makeBaselineConfigs()
	live := makeBaselineConfigs()
	live[0].Env["PORT"] = "9090"

	diffs := drift.CompareToBaseline(savedBaseline(base), live)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Field != "env.PORT" {
		t.Errorf("expected env.PORT diff, got %q", diffs[0].Field)
	}
}

func TestCompareToBaseline_MissingEnvKey(t *testing.T) {
	base := makeBaselineConfigs()
	live := makeBaselineConfigs()
	delete(live[0].Env, "LOG")

	diffs := drift.CompareToBaseline(savedBaseline(base), live)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Live != "<missing>" {
		t.Errorf("expected <missing>, got %q", diffs[0].Live)
	}
}

func TestCompareToBaseline_ServiceNotInBaseline(t *testing.T) {
	base := makeBaselineConfigs()
	live := append(makeBaselineConfigs(), drift.ServiceConfig{
		Name:  "newservice",
		Image: "new:1.0",
	})

	diffs := drift.CompareToBaseline(savedBaseline(base), live)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Service != "newservice" || diffs[0].Field != "existence" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}
