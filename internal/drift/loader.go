package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ServiceConfig represents the desired configuration of a service.
type ServiceConfig struct {
	Name        string            `json:"name" yaml:"name"`
	Image       string            `json:"image" yaml:"image"`
	Environment map[string]string `json:"environment" yaml:"environment"`
	Replicas    int               `json:"replicas" yaml:"replicas"`
}

// LoadConfig reads a service config from a JSON or YAML file.
func LoadConfig(path string) (*ServiceConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg ServiceConfig
	switch ext := filepath.Ext(path); ext {
	case ".json":
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parsing JSON config %q: %w", path, err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parsing YAML config %q: %w", path, err)
		}
	default:
		return nil, fmt.Errorf("unsupported config format %q (use .json, .yaml, or .yml)", ext)
	}

	if cfg.Name == "" {
		return nil, fmt.Errorf("config %q: missing required field 'name'", path)
	}
	if cfg.Image == "" {
		return nil, fmt.Errorf("config %q: missing required field 'image'", path)
	}

	return &cfg, nil
}

// LoadConfigDir loads all service configs from a directory.
func LoadConfigDir(dir string) ([]*ServiceConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading config directory %q: %w", dir, err)
	}

	var configs []*ServiceConfig
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			continue
		}
		cfg, err := LoadConfig(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		configs = append(configs, cfg)
	}
	return configs, nil
}
