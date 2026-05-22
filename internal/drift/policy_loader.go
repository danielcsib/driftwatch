package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadPolicy reads a Policy from a JSON or YAML file.
func LoadPolicy(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("policy: read %s: %w", path, err)
	}

	var p Policy
	switch ext := filepath.Ext(path); ext {
	case ".json":
		if err := json.Unmarshal(data, &p); err != nil {
			return nil, fmt.Errorf("policy: parse JSON %s: %w", path, err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &p); err != nil {
			return nil, fmt.Errorf("policy: parse YAML %s: %w", path, err)
		}
	default:
		return nil, fmt.Errorf("policy: unsupported extension %q", ext)
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("policy: invalid rules in %s: %w", path, err)
	}
	return &p, nil
}

// ApplyPolicy filters and annotates drift results according to the policy.
// Results whose every drifted field resolves to PolicyActionIgnore are dropped.
// If any field resolves to PolicyActionFail the result is marked as failed.
func ApplyPolicy(results []DriftResult, p *Policy) []DriftResult {
	if p == nil {
		return results
	}
	out := make([]DriftResult, 0, len(results))
	for _, r := range results {
		if !r.Drifted {
			out = append(out, r)
			continue
		}
		action := p.Evaluate(r.Service, r.Field)
		switch action {
		case PolicyActionIgnore:
			continue
		case PolicyActionFail:
			r.Field = "[FAIL] " + r.Field
			out = append(out, r)
		default:
			out = append(out, r)
		}
	}
	return out
}
