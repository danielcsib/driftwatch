package drift

import (
	"testing"
	"time"
)

func driftedResults() []DriftResult {
	return []DriftResult{
		{Service: "api", Drifted: true, Details: "image changed"},
	}
}

func cleanResults() []DriftResult {
	return []DriftResult{
		{Service: "api", Drifted: false},
	}
}

func TestHistory_RecordAndLen(t *testing.T) {
	h := NewHistory(10)
	if h.Len() != 0 {
		t.Fatalf("expected 0, got %d", h.Len())
	}
	h.Record(cleanResults())
	if h.Len() != 1 {
		t.Fatalf("expected 1, got %d", h.Len())
	}
}

func TestHistory_Latest_Empty(t *testing.T) {
	h := NewHistory(5)
	_, ok := h.Latest()
	if ok {
		t.Fatal("expected false for empty history")
	}
}

func TestHistory_Latest_ReturnsMostRecent(t *testing.T) {
	h := NewHistory(5)
	h.Record(cleanResults())
	time.Sleep(time.Millisecond)
	h.Record(driftedResults())

	ev, ok := h.Latest()
	if !ok {
		t.Fatal("expected event")
	}
	if !ev.HasDrift {
		t.Error("expected latest event to have drift")
	}
}

func TestHistory_HasDrift_Flag(t *testing.T) {
	h := NewHistory(10)
	h.Record(driftedResults())

	ev, _ := h.Latest()
	if !ev.HasDrift {
		t.Error("expected HasDrift=true")
	}

	h.Record(cleanResults())
	ev, _ = h.Latest()
	if ev.HasDrift {
		t.Error("expected HasDrift=false")
	}
}

func TestHistory_BoundedBuffer(t *testing.T) {
	h := NewHistory(3)
	for i := 0; i < 5; i++ {
		h.Record(cleanResults())
	}
	if h.Len() != 3 {
		t.Fatalf("expected 3 (maxLen), got %d", h.Len())
	}
}

func TestHistory_All_ReturnsCopy(t *testing.T) {
	h := NewHistory(10)
	h.Record(cleanResults())
	h.Record(driftedResults())

	all := h.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 events, got %d", len(all))
	}
	// Mutating the returned slice must not affect internal state.
	all[0] = DriftEvent{}
	if h.Len() != 2 {
		t.Error("internal history was mutated")
	}
}

func TestNewHistory_DefaultMaxLen(t *testing.T) {
	h := NewHistory(0)
	if h.maxLen != 100 {
		t.Errorf("expected default maxLen 100, got %d", h.maxLen)
	}
}
