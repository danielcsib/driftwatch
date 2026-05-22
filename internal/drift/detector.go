package drift

import (
	"crypto/sha256"
	"fmt"
	"sort"
)

// ServiceConfig represents the desired and live configuration for a service.
type ServiceConfig struct {
	Name     string            `json:"name"     yaml:"name"`
	Image    string            `json:"image"    yaml:"image"`
	Replicas int               `json:"replicas" yaml:"replicas"`
	Env      map[string]string `json:"env"      yaml:"env"`
}

// DetectResult holds the outcome of comparing desired vs live config.
type DetectResult struct {
	Service       string
	Drifted       bool
	DriftedFields []string
}

// Detector compares desired configs against live configs.
type Detector struct{}

// NewDetector returns a new Detector.
func NewDetector() *Detector { return &Detector{} }

// Detect compares desired to live and returns a DetectResult.
func (d *Detector) Detect(desired, live ServiceConfig) DetectResult {
	result := DetectResult{Service: desired.Name}
	var fields []string

	if desired.Image != live.Image {
		fields = append(fields, "image")
	}
	if desired.Replicas != live.Replicas {
		fields = append(fields, "replicas")
	}
	for k, dv := range desired.Env {
		lv, ok := live.Env[k]
		if !ok {
			fields = append(fields, "env."+k+" (missing)")
		} else if dv != lv {
			fields = append(fields, "env."+k)
		}
	}
	sort.Strings(fields)
	result.DriftedFields = fields
	result.Drifted = len(fields) > 0
	return result
}

// Checksum returns a SHA-256 hex digest of the canonical config representation.
func Checksum(cfg ServiceConfig) string {
	keys := make([]string, 0, len(cfg.Env))
	for k := range cfg.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	fmt.Fprintf(h, "%s|%s|%d", cfg.Name, cfg.Image, cfg.Replicas)
	for _, k := range keys {
		fmt.Fprintf(h, "|%s=%s", k, cfg.Env[k])
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ConfigsMatch returns true when desired and live produce the same checksum.
func ConfigsMatch(desired, live ServiceConfig) bool {
	return Checksum(desired) == Checksum(live)
}
