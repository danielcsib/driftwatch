package drift

import (
	"fmt"
	"strings"
	"time"
)

// Summary aggregates drift detection results into a human-readable report.
type Summary struct {
	GeneratedAt  time.Time
	TotalChecked int
	DriftedCount int
	CleanCount   int
	Services     []ServiceSummary
}

// ServiceSummary holds drift info for a single service.
type ServiceSummary struct {
	Name    string
	Drifted bool
	Fields  []string
}

// NewSummary builds a Summary from a slice of DetectResult.
func NewSummary(results []DetectResult) Summary {
	s := Summary{
		GeneratedAt:  time.Now().UTC(),
		TotalChecked: len(results),
	}
	for _, r := range results {
		svc := ServiceSummary{
			Name:    r.Service,
			Drifted: r.Drifted,
			Fields:  r.DriftedFields,
		}
		if r.Drifted {
			s.DriftedCount++
		} else {
			s.CleanCount++
		}
		s.Services = append(s.Services, svc)
	}
	return s
}

// String returns a compact text representation of the Summary.
func (s Summary) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "DriftWatch Summary [%s]\n", s.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "  Checked: %d | Drifted: %d | Clean: %d\n", s.TotalChecked, s.DriftedCount, s.CleanCount)
	for _, svc := range s.Services {
		if svc.Drifted {
			fmt.Fprintf(&b, "  [DRIFT] %s — fields: %s\n", svc.Name, strings.Join(svc.Fields, ", "))
		} else {
			fmt.Fprintf(&b, "  [OK]    %s\n", svc.Name)
		}
	}
	return b.String()
}
