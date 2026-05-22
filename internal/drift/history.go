package drift

import (
	"sync"
	"time"
)

// DriftEvent records a single drift detection result at a point in time.
type DriftEvent struct {
	Timestamp time.Time
	Results   []DriftResult
	HasDrift  bool
}

// History stores a bounded ring-buffer of past drift events.
type History struct {
	mu     sync.RWMutex
	events []DriftEvent
	maxLen int
}

// NewHistory creates a History that retains at most maxLen events.
// If maxLen <= 0 it defaults to 100.
func NewHistory(maxLen int) *History {
	if maxLen <= 0 {
		maxLen = 100
	}
	return &History{maxLen: maxLen}
}

// Record appends a new event, evicting the oldest if the buffer is full.
func (h *History) Record(results []DriftResult) {
	h.mu.Lock()
	defer h.mu.Unlock()

	hasDrift := false
	for _, r := range results {
		if r.Drifted {
			hasDrift = true
			break
		}
	}

	event := DriftEvent{
		Timestamp: time.Now().UTC(),
		Results:   results,
		HasDrift:  hasDrift,
	}

	if len(h.events) >= h.maxLen {
		h.events = h.events[1:]
	}
	h.events = append(h.events, event)
}

// All returns a copy of all stored events, oldest first.
func (h *History) All() []DriftEvent {
	h.mu.RLock()
	defer h.mu.RUnlock()

	copy := make([]DriftEvent, len(h.events))
	for i, e := range h.events {
		copy[i] = e
	}
	return copy
}

// Latest returns the most recent event and true, or a zero value and false
// if no events have been recorded yet.
func (h *History) Latest() (DriftEvent, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.events) == 0 {
		return DriftEvent{}, false
	}
	return h.events[len(h.events)-1], true
}

// Len returns the number of events currently stored.
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.events)
}
