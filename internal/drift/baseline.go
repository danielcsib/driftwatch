package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Baseline represents a saved reference state for a set of service configs.
type Baseline struct {
	CreatedAt time.Time        `json:"created_at"`
	Label     string           `json:"label"`
	Configs   []ServiceConfig  `json:"configs"`
}

// BaselineStore persists and retrieves baselines on disk.
type BaselineStore struct {
	dir string
}

// NewBaselineStore creates a BaselineStore rooted at dir.
func NewBaselineStore(dir string) (*BaselineStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("baseline: create dir: %w", err)
	}
	return &BaselineStore{dir: dir}, nil
}

// Save writes a baseline to disk using label as the filename key.
func (s *BaselineStore) Save(label string, configs []ServiceConfig) error {
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Configs:   configs,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	path := filepath.Join(s.dir, label+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write: %w", err)
	}
	return nil
}

// Load retrieves a baseline by label.
func (s *BaselineStore) Load(label string) (*Baseline, error) {
	path := filepath.Join(s.dir, label+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read %q: %w", label, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// List returns all baseline labels stored in the directory.
func (s *BaselineStore) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("baseline: list: %w", err)
	}
	var labels []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			labels = append(labels, e.Name()[:len(e.Name())-5])
		}
	}
	return labels, nil
}
