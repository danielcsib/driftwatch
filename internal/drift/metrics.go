package drift

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Metrics tracks runtime statistics for drift checks.
type Metrics struct {
	mu           sync.Mutex
	TotalChecks  int
	DriftedCount int
	CleanCount   int
	LastCheckAt  time.Time
	LastDriftAt  time.Time
	Errors       int
	out          io.Writer
}

// NewMetrics creates a Metrics instance writing summaries to out.
// If out is nil, os.Stderr is used.
func NewMetrics(out io.Writer) *Metrics {
	if out == nil {
		out = os.Stderr
	}
	return &Metrics{out: out}
}

// Record updates metrics based on a slice of DriftResult.
func (m *Metrics) Record(results []DriftResult) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalChecks++
	m.LastCheckAt = time.Now()

	drifted := 0
	for _, r := range results {
		if r.Drifted {
			drifted++
		}
	}

	if drifted > 0 {
		m.DriftedCount += drifted
		m.LastDriftAt = time.Now()
	} else {
		m.CleanCount++
	}
}

// RecordError increments the error counter.
func (m *Metrics) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Errors++
}

// Print writes a human-readable metrics summary to the configured writer.
func (m *Metrics) Print() {
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Fprintf(m.out, "[driftwatch metrics] checks=%d drifted_services=%d clean_checks=%d errors=%d\n",
		m.TotalChecks, m.DriftedCount, m.CleanCount, m.Errors)

	if !m.LastCheckAt.IsZero() {
		fmt.Fprintf(m.out, "  last_check=%s\n", m.LastCheckAt.Format(time.RFC3339))
	}
	if !m.LastDriftAt.IsZero() {
		fmt.Fprintf(m.out, "  last_drift=%s\n", m.LastDriftAt.Format(time.RFC3339))
	}
}

// Snapshot returns an immutable copy of the current metrics.
func (m *Metrics) Snapshot() Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	return Metrics{
		TotalChecks:  m.TotalChecks,
		DriftedCount: m.DriftedCount,
		CleanCount:   m.CleanCount,
		LastCheckAt:  m.LastCheckAt,
		LastDriftAt:  m.LastDriftAt,
		Errors:       m.Errors,
	}
}
