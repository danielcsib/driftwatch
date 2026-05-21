package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestLoadConfig_JSON(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, "svc.json", `{"name":"api","image":"api:v1","replicas":2,"environment":{"PORT":"8080"}}`)

	cfg, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "api" || cfg.Image != "api:v1" || cfg.Replicas != 2 {
		t.Errorf("unexpected config: %+v", cfg)
	}
	if cfg.Environment["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", cfg.Environment["PORT"])
	}
}

func TestLoadConfig_YAML(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, "svc.yaml", "name: worker\nimage: worker:latest\nreplicas: 1\nenvironment:\n  LOG_LEVEL: debug\n")

	cfg, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "worker" || cfg.Image != "worker:latest" {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestLoadConfig_MissingName(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, "bad.json", `{"image":"api:v1"}`)

	_, err := LoadConfig(p)
	if err == nil {
		t.Fatal("expected error for missing name, got nil")
	}
}

func TestLoadConfig_UnsupportedExt(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, "svc.toml", "name = \"api\"")

	_, err := LoadConfig(p)
	if err == nil {
		t.Fatal("expected error for unsupported extension")
	}
}

func TestLoadConfigDir(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "a.json", `{"name":"a","image":"a:1"}`)
	writeFile(t, dir, "b.yaml", "name: b\nimage: b:1\n")
	writeFile(t, dir, "ignore.txt", "not a config")

	configs, err := LoadConfigDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(configs) != 2 {
		t.Errorf("expected 2 configs, got %d", len(configs))
	}
}
