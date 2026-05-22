package drift

import (
	"fmt"
	"io"
	"os"
	"time"
)

// AlertLevel represents the severity of a drift alert.
type AlertLevel string

const (
	AlertLevelWarn  AlertLevel = "WARN"
	AlertLevelError AlertLevel = "ERROR"
)

// Alert represents a drift alert emitted when drift is detected.
type Alert struct {
	Timestamp   time.Time  `json:"timestamp"`
	ServiceName string     `json:"service_name"`
	Level       AlertLevel `json:"level"`
	Message     string     `json:"message"`
	Drifts      []DriftDetail `json:"drifts"`
}

// DriftDetail describes a single field that has drifted.
type DriftDetail struct {
	Field    string `json:"field"`
	Wanted   string `json:"wanted"`
	Actual   string `json:"actual"`
}

// Alerter sends alerts when drift is detected.
type Alerter struct {
	out   io.Writer
	level AlertLevel
}

// NewAlerter creates an Alerter writing to out at the given level.
// If out is nil, os.Stderr is used.
func NewAlerter(out io.Writer, level AlertLevel) *Alerter {
	if out == nil {
		out = os.Stderr
	}
	return &Alerter{out: out, level: level}
}

// Send writes a formatted alert to the configured writer.
func (a *Alerter) Send(alert Alert) error {
	_, err := fmt.Fprintf(
		a.out,
		"[%s] %s | service=%s drifts=%d\n",
		alert.Level,
		alert.Timestamp.Format(time.RFC3339),
		alert.ServiceName,
		len(alert.Drifts),
	)
	return err
}

// BuildAlert constructs an Alert from a Result produced by the Detector.
func BuildAlert(result Result, level AlertLevel) Alert {
	details := make([]DriftDetail, 0, len(result.Differences))
	for _, d := range result.Differences {
		details = append(details, DriftDetail{
			Field:  d.Field,
			Wanted: d.Wanted,
			Actual: d.Actual,
		})
	}
	return Alert{
		Timestamp:   time.Now().UTC(),
		ServiceName: result.ServiceName,
		Level:       level,
		Message:     fmt.Sprintf("drift detected in service %q", result.ServiceName),
		Drifts:      details,
	}
}
