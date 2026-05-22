package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPolicy_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	content := `{"rules":[{"service":"api","field":"image","action":"ignore"}]}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := LoadPolicy(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(p.Rules))
	}
	if p.Rules[0].Action != PolicyActionIgnore {
		t.Fatalf("expected ignore, got %q", p.Rules[0].Action)
	}
}

func TestLoadPolicy_YAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.yaml")
	content := "rules:\n  - service: api\n    field: image\n    action: fail\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := LoadPolicy(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Rules[0].Action != PolicyActionFail {
		t.Fatalf("expected fail, got %q", p.Rules[0].Action)
	}
}

func TestLoadPolicy_UnsupportedExt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.toml")
	_ = os.WriteFile(path, []byte("x"), 0o644)
	if _, err := LoadPolicy(path); err == nil {
		t.Fatal("expected error for unsupported extension")
	}
}

func TestLoadPolicy_InvalidAction(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	content := `{"rules":[{"service":"api","field":"image","action":"destroy"}]}`
	_ = os.WriteFile(path, []byte(content), 0o644)
	if _, err := LoadPolicy(path); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestLoadPolicy_MissingFile(t *testing.T) {
	if _, err := LoadPolicy("/nonexistent/policy.json"); err == nil {
		t.Fatal("expected error for missing file")
	}
}
