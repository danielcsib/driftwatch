package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of drift detection results.
type Snapshot struct {
	Timestamp time.Time       `json:"timestamp"`
	Results   []DetectResult  `json:"results"`
}

// SnapshotStore persists and retrieves drift snapshots to/from disk.
type SnapshotStore struct {
	dir string
}

// NewSnapshotStore creates a SnapshotStore that writes to the given directory.
func NewSnapshotStore(dir string) (*SnapshotStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &SnapshotStore{dir: dir}, nil
}

// Save writes a snapshot to disk as a JSON file named by its timestamp.
func (s *SnapshotStore) Save(snap Snapshot) error {
	filename := snap.Timestamp.UTC().Format("20060102T150405Z") + ".json"
	path := filepath.Join(s.dir, filename)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Latest returns the most recently saved snapshot, or an error if none exist.
func (s *SnapshotStore) Latest() (Snapshot, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read dir: %w", err)
	}

	var last os.DirEntry
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			last = e
		}
	}
	if last == nil {
		return Snapshot{}, fmt.Errorf("snapshot: no snapshots found")
	}

	data, err := os.ReadFile(filepath.Join(s.dir, last.Name()))
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read file: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: decode: %w", err)
	}
	return snap, nil
}

// List returns all stored snapshots in chronological order.
func (s *SnapshotStore) List() ([]Snapshot, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read dir: %w", err)
	}

	var snaps []Snapshot
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("snapshot: read %s: %w", e.Name(), err)
		}
		var snap Snapshot
		if err := json.Unmarshal(data, &snap); err != nil {
			return nil, fmt.Errorf("snapshot: decode %s: %w", e.Name(), err)
		}
		snaps = append(snaps, snap)
	}
	return snaps, nil
}
