package drift_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

func TestAlerter_Send_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	a := drift.NewAlerter(&buf, drift.AlertLevelWarn)

	alert := drift.Alert{
		Timestamp:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		ServiceName: "api-gateway",
		Level:       drift.AlertLevelWarn,
		Message:     "drift detected",
		Drifts: []drift.DriftDetail{
			{Field: "image", Wanted: "v1.2", Actual: "v1.1"},
		},
	}

	if err := a.Send(alert); err != nil {
		t.Fatalf("Send() unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %s", out)
	}
	if !strings.Contains(out, "api-gateway") {
		t.Errorf("expected service name in output, got: %s", out)
	}
	if !strings.Contains(out, "drifts=1") {
		t.Errorf("expected drifts=1 in output, got: %s", out)
	}
}

func TestBuildAlert_PopulatesFields(t *testing.T) {
	result := drift.Result{
		ServiceName: "worker",
		HasDrift:    true,
		Differences: []drift.Difference{
			{Field: "replicas", Wanted: "3", Actual: "1"},
			{Field: "image", Wanted: "v2.0", Actual: "v1.9"},
		},
	}

	alert := drift.BuildAlert(result, drift.AlertLevelError)

	if alert.ServiceName != "worker" {
		t.Errorf("expected service name 'worker', got %q", alert.ServiceName)
	}
	if alert.Level != drift.AlertLevelError {
		t.Errorf("expected level ERROR, got %q", alert.Level)
	}
	if len(alert.Drifts) != 2 {
		t.Errorf("expected 2 drifts, got %d", len(alert.Drifts))
	}
	if alert.Drifts[0].Field != "replicas" {
		t.Errorf("expected first drift field 'replicas', got %q", alert.Drifts[0].Field)
	}
	if alert.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestNewAlerter_DefaultsToStderr(t *testing.T) {
	a := drift.NewAlerter(nil, drift.AlertLevelWarn)
	if a == nil {
		t.Fatal("expected non-nil Alerter")
	}
}
