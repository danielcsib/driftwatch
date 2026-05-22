package drift

import (
	"strings"
	"testing"
)

func driftedSummaryResults() []DetectResult {
	return []DetectResult{
		{Service: "api", Drifted: true, DriftedFields: []string{"image", "env.PORT"}},
		{Service: "worker", Drifted: false, DriftedFields: nil},
		{Service: "cache", Drifted: true, DriftedFields: []string{"replicas"}},
	}
}

func TestNewSummary_Counts(t *testing.T) {
	results := driftedSummaryResults()
	s := NewSummary(results)

	if s.TotalChecked != 3 {
		t.Errorf("TotalChecked: want 3, got %d", s.TotalChecked)
	}
	if s.DriftedCount != 2 {
		t.Errorf("DriftedCount: want 2, got %d", s.DriftedCount)
	}
	if s.CleanCount != 1 {
		t.Errorf("CleanCount: want 1, got %d", s.CleanCount)
	}
}

func TestNewSummary_Services(t *testing.T) {
	s := NewSummary(driftedSummaryResults())

	if len(s.Services) != 3 {
		t.Fatalf("expected 3 services, got %d", len(s.Services))
	}
	if !s.Services[0].Drifted {
		t.Error("api should be drifted")
	}
	if s.Services[1].Drifted {
		t.Error("worker should be clean")
	}
}

func TestNewSummary_Empty(t *testing.T) {
	s := NewSummary(nil)
	if s.TotalChecked != 0 || s.DriftedCount != 0 || s.CleanCount != 0 {
		t.Error("empty results should yield zero counts")
	}
}

func TestSummary_String_ContainsDrift(t *testing.T) {
	s := NewSummary(driftedSummaryResults())
	out := s.String()

	if !strings.Contains(out, "[DRIFT] api") {
		t.Error("expected [DRIFT] api in output")
	}
	if !strings.Contains(out, "[OK]    worker") {
		t.Error("expected [OK] worker in output")
	}
	if !strings.Contains(out, "image") {
		t.Error("expected drifted field 'image' in output")
	}
}

func TestSummary_String_Header(t *testing.T) {
	s := NewSummary(driftedSummaryResults())
	out := s.String()

	if !strings.Contains(out, "DriftWatch Summary") {
		t.Error("expected header in output")
	}
	if !strings.Contains(out, "Checked: 3") {
		t.Error("expected checked count in output")
	}
}
