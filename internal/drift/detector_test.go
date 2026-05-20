package drift

import (
	"testing"
)

func baseDesired() *ServiceConfig {
	return &ServiceConfig{
		Name:  "api",
		Image: "myrepo/api:v1.2.3",
		Env:   map[string]string{"PORT": "8080", "LOG_LEVEL": "info"},
		Labels: map[string]string{"team": "platform"},
	}
}

func TestDetect_NoDrift(t *testing.T) {
	d := NewDetector()
	desired := baseDesired()
	actual := baseDesired()

	result := d.Detect(desired, actual)

	if result.HasDrift {
		t.Errorf("expected no drift, got diffs: %v", result.Diffs)
	}
}

func TestDetect_ImageDrift(t *testing.T) {
	d := NewDetector()
	desired := baseDesired()
	actual := baseDesired()
	actual.Image = "myrepo/api:v1.2.0"

	result := d.Detect(desired, actual)

	if !result.HasDrift {
		t.Fatal("expected drift but none reported")
	}
	if len(result.Diffs) != 1 || result.Diffs[0] != `image: want "myrepo/api:v1.2.3", got "myrepo/api:v1.2.0"` {
		t.Errorf("unexpected diffs: %v", result.Diffs)
	}
}

func TestDetect_EnvDrift(t *testing.T) {
	d := NewDetector()
	desired := baseDesired()
	actual := baseDesired()
	actual.Env["LOG_LEVEL"] = "debug"

	result := d.Detect(desired, actual)

	if !result.HasDrift {
		t.Fatal("expected drift but none reported")
	}
}

func TestDetect_MissingEnvKey(t *testing.T) {
	d := NewDetector()
	desired := baseDesired()
	actual := baseDesired()
	delete(actual.Env, "PORT")

	result := d.Detect(desired, actual)

	if !result.HasDrift {
		t.Fatal("expected drift for missing env key")
	}
}

func TestChecksum_Consistency(t *testing.T) {
	cfg := baseDesired()
	h1, err := Checksum(cfg)
	if err != nil {
		t.Fatalf("checksum error: %v", err)
	}
	h2, err := Checksum(cfg)
	if err != nil {
		t.Fatalf("checksum error: %v", err)
	}
	if h1 != h2 {
		t.Errorf("checksums differ for identical configs: %s vs %s", h1, h2)
	}
}

func TestChecksum_DifferentConfigs(t *testing.T) {
	c1 := baseDesired()
	c2 := baseDesired()
	c2.Image = "myrepo/api:v2.0.0"

	h1, _ := Checksum(c1)
	h2, _ := Checksum(c2)

	if h1 == h2 {
		t.Error("expected different checksums for different configs")
	}
}
