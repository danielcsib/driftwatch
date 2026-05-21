package drift

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// ServiceConfig represents the expected configuration for a service.
type ServiceConfig struct {
	Name   string            `json:"name"`
	Image  string            `json:"image"`
	Env    map[string]string `json:"env"`
	Labels map[string]string `json:"labels"`
}

// DriftResult holds the comparison result between desired and actual state.
type DriftResult struct {
	ServiceName string
	HasDrift    bool
	Diffs       []string
}

// Detector compares deployed service state against source definitions.
type Detector struct{}

// NewDetector creates a new Detector instance.
func NewDetector() *Detector {
	return &Detector{}
}

// Detect compares the desired config against the actual config and returns drift details.
// It checks for differences in image, environment variables, and labels.
// Keys present in actual but absent in desired are not flagged as drift;
// only keys required by desired are enforced.
func (d *Detector) Detect(desired, actual *ServiceConfig) DriftResult {
	result := DriftResult{
		ServiceName: desired.Name,
		HasDrift:    false,
		Diffs:       []string{},
	}

	if desired.Image != actual.Image {
		result.HasDrift = true
		result.Diffs = append(result.Diffs, fmt.Sprintf("image: want %q, got %q", desired.Image, actual.Image))
	}

	for k, wantVal := range desired.Env {
		if gotVal, ok := actual.Env[k]; !ok {
			result.HasDrift = true
			result.Diffs = append(result.Diffs, fmt.Sprintf("env[%s]: missing (want %q)", k, wantVal))
		} else if gotVal != wantVal {
			result.HasDrift = true
			result.Diffs = append(result.Diffs, fmt.Sprintf("env[%s]: want %q, got %q", k, wantVal, gotVal))
		}
	}

	for k, wantVal := range desired.Labels {
		if gotVal, ok := actual.Labels[k]; !ok {
			result.HasDrift = true
			result.Diffs = append(result.Diffs, fmt.Sprintf("label[%s]: missing (want %q)", k, wantVal))
		} else if gotVal != wantVal {
			result.HasDrift = true
			result.Diffs = append(result.Diffs, fmt.Sprintf("label[%s]: want %q, got %q", k, wantVal, gotVal))
		}
	}

	return result
}

// Checksum returns a SHA-256 hash of the service config for quick equality checks.
func Checksum(cfg *ServiceConfig) (string, error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("checksum marshal: %w", err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

// ConfigsMatch returns true if the two service configs produce identical checksums.
// It returns an error if either config cannot be marshalled.
func ConfigsMatch(a, b *ServiceConfig) (bool, error) {
	sumA, err := Checksum(a)
	if err != nil {
		return false, fmt.Errorf("configs match: %w", err)
	}
	sumB, err := Checksum(b)
	if err != nil {
		return false, fmt.Errorf("configs match: %w", err)
	}
	return sumA == sumB, nil
}
