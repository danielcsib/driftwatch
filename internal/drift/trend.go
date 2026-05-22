package drift

import (
	"fmt"
	"time"
)

// TrendDirection indicates whether drift is increasing, decreasing, or stable.
type TrendDirection string

const (
	TrendIncreasing TrendDirection = "increasing"
	TrendDecreasing TrendDirection = "decreasing"
	TrendStable     TrendDirection = "stable"
)

// TrendPoint represents the drift count at a specific point in time.
type TrendPoint struct {
	Timestamp  time.Time
	DriftCount int
}

// Trend summarises drift movement across a series of history snapshots.
type Trend struct {
	Points    []TrendPoint
	Direction TrendDirection
	Delta     int // difference between first and last point
}

// String returns a human-readable summary of the trend.
func (t Trend) String() string {
	if len(t.Points) == 0 {
		return "trend: no data"
	}
	return fmt.Sprintf("trend: %s (delta %+d over %d checks)",
		t.Direction, t.Delta, len(t.Points))
}

// AnalyseTrend builds a Trend from a slice of HistoryEntry values.
// Each entry contributes one TrendPoint based on how many results drifted.
func AnalyseTrend(entries []HistoryEntry) Trend {
	if len(entries) == 0 {
		return Trend{Direction: TrendStable}
	}

	points := make([]TrendPoint, 0, len(entries))
	for _, e := range entries {
		count := 0
		for _, r := range e.Results {
			if r.Drifted {
				count++
			}
		}
		points = append(points, TrendPoint{
			Timestamp:  e.CheckedAt,
			DriftCount: count,
		})
	}

	first := points[0].DriftCount
	last := points[len(points)-1].DriftCount
	delta := last - first

	var dir TrendDirection
	switch {
	case delta > 0:
		dir = TrendIncreasing
	case delta < 0:
		dir = TrendDecreasing
	default:
		dir = TrendStable
	}

	return Trend{
		Points:    points,
		Direction: dir,
		Delta:     delta,
	}
}
