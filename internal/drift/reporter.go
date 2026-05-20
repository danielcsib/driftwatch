package drift

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// DriftResult holds the outcome of a drift check for a single service.
type DriftResult struct {
	Service   string        `json:"service"`
	Drifted   bool          `json:"drifted"`
	Diffs     []Diff        `json:"diffs,omitempty"`
	CheckedAt time.Time     `json:"checked_at"`
}

// Diff describes a single field that has drifted.
type Diff struct {
	Field    string `json:"field"`
	Wanted   string `json:"wanted"`
	Actual   string `json:"actual"`
}

// Reporter formats and writes drift results to an io.Writer.
type Reporter struct {
	w      io.Writer
	format string // "text" or "json"
}

// NewReporter creates a Reporter that writes to w using the given format.
// Supported formats: "text", "json". Defaults to "text" for unknown values.
func NewReporter(w io.Writer, format string) *Reporter {
	return &Reporter{w: w, format: format}
}

// Report writes the drift result to the underlying writer.
func (r *Reporter) Report(result DriftResult) error {
	switch r.format {
	case "json":
		return r.reportJSON(result)
	default:
		return r.reportText(result)
	}
}

func (r *Reporter) reportText(result DriftResult) error {
	if !result.Drifted {
		_, err := fmt.Fprintf(r.w, "[%s] %s: OK (no drift)\n",
			result.CheckedAt.Format(time.RFC3339), result.Service)
		return err
	}
	_, err := fmt.Fprintf(r.w, "[%s] %s: DRIFTED (%d change(s))\n",
		result.CheckedAt.Format(time.RFC3339), result.Service, len(result.Diffs))
	if err != nil {
		return err
	}
	for _, d := range result.Diffs {
		_, err = fmt.Fprintf(r.w, "  field=%s wanted=%q actual=%q\n", d.Field, d.Wanted, d.Actual)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reporter) reportJSON(result DriftResult) error {
	enc := json.NewEncoder(r.w)
	return enc.Encode(result)
}
