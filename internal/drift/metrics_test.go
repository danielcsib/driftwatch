package drift

import (
	"bytes"
	"strings"
	"testing"
)

func driftedMetricResults() []DriftResult {
	return []DriftResult{
		{Service: "api", Drifted: true},
		{Service: "worker", Drifted: false},
	}
}

func cleanMetricResults() []DriftResult {
	return []DriftResult{
		{Service: "api", Drifted: false},
		{Service: "worker", Drifted: false},
	}
}

func TestMetrics_RecordDrifted(t *testing.T) {
	m := NewMetrics(nil)
	m.Record(driftedMetricResults())

	snap := m.Snapshot()
	if snap.TotalChecks != 1 {
		t.Errorf("expected TotalChecks=1, got %d", snap.TotalChecks)
	}
	if snap.DriftedCount != 1 {
		t.Errorf("expected DriftedCount=1, got %d", snap.DriftedCount)
	}
	if snap.CleanCount != 0 {
		t.Errorf("expected CleanCount=0, got %d", snap.CleanCount)
	}
	if snap.LastDriftAt.IsZero() {
		t.Error("expected LastDriftAt to be set")
	}
}

func TestMetrics_RecordClean(t *testing.T) {
	m := NewMetrics(nil)
	m.Record(cleanMetricResults())

	snap := m.Snapshot()
	if snap.CleanCount != 1 {
		t.Errorf("expected CleanCount=1, got %d", snap.CleanCount)
	}
	if snap.DriftedCount != 0 {
		t.Errorf("expected DriftedCount=0, got %d", snap.DriftedCount)
	}
	if !snap.LastDriftAt.IsZero() {
		t.Error("expected LastDriftAt to be zero on clean run")
	}
}

func TestMetrics_RecordError(t *testing.T) {
	m := NewMetrics(nil)
	m.RecordError()
	m.RecordError()

	snap := m.Snapshot()
	if snap.Errors != 2 {
		t.Errorf("expected Errors=2, got %d", snap.Errors)
	}
}

func TestMetrics_Print(t *testing.T) {
	var buf bytes.Buffer
	m := NewMetrics(&buf)
	m.Record(driftedMetricResults())
	m.RecordError()
	m.Print()

	out := buf.String()
	if !strings.Contains(out, "checks=1") {
		t.Errorf("expected checks=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "drifted_services=1") {
		t.Errorf("expected drifted_services=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "errors=1") {
		t.Errorf("expected errors=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "last_drift=") {
		t.Errorf("expected last_drift= in output, got: %s", out)
	}
}

func TestNewMetrics_DefaultsToStderr(t *testing.T) {
	m := NewMetrics(nil)
	if m.out == nil {
		t.Error("expected non-nil writer when nil passed to NewMetrics")
	}
}
