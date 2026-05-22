package drift

import (
	"testing"
	"time"
)

func makeEntry(t time.Time, driftedCount, cleanCount int) HistoryEntry {
	results := make([]DriftResult, 0, driftedCount+cleanCount)
	for i := 0; i < driftedCount; i++ {
		results = append(results, DriftResult{Service: "svc", Drifted: true})
	}
	for i := 0; i < cleanCount; i++ {
		results = append(results, DriftResult{Service: "svc", Drifted: false})
	}
	return HistoryEntry{CheckedAt: t, Results: results}
}

func TestAnalyseTrend_Empty(t *testing.T) {
	trend := AnalyseTrend(nil)
	if trend.Direction != TrendStable {
		t.Errorf("expected stable, got %s", trend.Direction)
	}
	if len(trend.Points) != 0 {
		t.Errorf("expected no points, got %d", len(trend.Points))
	}
}

func TestAnalyseTrend_Stable(t *testing.T) {
	now := time.Now()
	entries := []HistoryEntry{
		makeEntry(now.Add(-2*time.Minute), 2, 1),
		makeEntry(now.Add(-1*time.Minute), 2, 1),
		makeEntry(now, 2, 1),
	}
	trend := AnalyseTrend(entries)
	if trend.Direction != TrendStable {
		t.Errorf("expected stable, got %s", trend.Direction)
	}
	if trend.Delta != 0 {
		t.Errorf("expected delta 0, got %d", trend.Delta)
	}
	if len(trend.Points) != 3 {
		t.Errorf("expected 3 points, got %d", len(trend.Points))
	}
}

func TestAnalyseTrend_Increasing(t *testing.T) {
	now := time.Now()
	entries := []HistoryEntry{
		makeEntry(now.Add(-2*time.Minute), 1, 2),
		makeEntry(now.Add(-1*time.Minute), 2, 1),
		makeEntry(now, 4, 0),
	}
	trend := AnalyseTrend(entries)
	if trend.Direction != TrendIncreasing {
		t.Errorf("expected increasing, got %s", trend.Direction)
	}
	if trend.Delta != 3 {
		t.Errorf("expected delta 3, got %d", trend.Delta)
	}
}

func TestAnalyseTrend_Decreasing(t *testing.T) {
	now := time.Now()
	entries := []HistoryEntry{
		makeEntry(now.Add(-2*time.Minute), 5, 0),
		makeEntry(now.Add(-1*time.Minute), 3, 2),
		makeEntry(now, 1, 4),
	}
	trend := AnalyseTrend(entries)
	if trend.Direction != TrendDecreasing {
		t.Errorf("expected decreasing, got %s", trend.Direction)
	}
	if trend.Delta != -4 {
		t.Errorf("expected delta -4, got %d", trend.Delta)
	}
}

func TestTrend_String(t *testing.T) {
	now := time.Now()
	entries := []HistoryEntry{
		makeEntry(now.Add(-time.Minute), 1, 0),
		makeEntry(now, 3, 0),
	}
	trend := AnalyseTrend(entries)
	s := trend.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestTrend_String_NoData(t *testing.T) {
	trend := Trend{}
	if trend.String() != "trend: no data" {
		t.Errorf("unexpected string: %s", trend.String())
	}
}
