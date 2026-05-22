package drift

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestReporter_TextNoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, "text")
	result := DriftResult{
		Service:   "api",
		Drifted:   false,
		CheckedAt: fixedTime,
	}
	if err := r.Report(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "OK (no drift)") {
		t.Errorf("expected 'OK (no drift)' in output, got: %q", got)
	}
	if !strings.Contains(got, "api") {
		t.Errorf("expected service name 'api' in output, got: %q", got)
	}
}

func TestReporter_TextDrifted(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, "text")
	result := DriftResult{
		Service: "worker",
		Drifted: true,
		Diffs: []Diff{
			{Field: "image", Wanted: "app:v2", Actual: "app:v1"},
			{Field: "env.LOG_LEVEL", Wanted: "info", Actual: "debug"},
		},
		CheckedAt: fixedTime,
	}
	if err := r.Report(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "DRIFTED") {
		t.Errorf("expected 'DRIFTED' in output, got: %q", got)
	}
	if !strings.Contains(got, "image") {
		t.Errorf("expected field 'image' in output, got: %q", got)
	}
	if !strings.Contains(got, "2 change(s)") {
		t.Errorf("expected '2 change(s)' in output, got: %q", got)
	}
}

func TestReporter_JSONOutput(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, "json")
	result := DriftResult{
		Service: "gateway",
		Drifted: true,
		Diffs:   []Diff{{Field: "replicas", Wanted: "3", Actual: "1"}},
		CheckedAt: fixedTime,
	}
	if err := r.Report(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded DriftResult
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to decode JSON output: %v", err)
	}
	if decoded.Service != "gateway" {
		t.Errorf("expected service 'gateway', got %q", decoded.Service)
	}
	if len(decoded.Diffs) != 1 {
		t.Errorf("expected 1 diff, got %d", len(decoded.Diffs))
	}
}

func TestReporter_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, "xml")
	result := DriftResult{
		Service:   "api",
		Drifted:   false,
		CheckedAt: fixedTime,
	}
	if err := r.Report(result); err == nil {
		t.Error("expected error for unknown format 'xml', got nil")
	}
}
