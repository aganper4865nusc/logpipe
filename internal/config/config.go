// Package config loads and validates logpipe configuration from a YAML file.
package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const defaultWorkers = 4

// RateLimitConfig controls per-source throughput.
type RateLimitConfig struct {
	Rate  float64 `yaml:"rate"`  // lines per second (0 = unlimited)
	Burst float64 `yaml:"burst"` // burst capacity
}

// SourceConfig describes a single tailed file source.
type SourceConfig struct {
	Name      string          `yaml:"name"`
	Path      string          `yaml:"path"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
}

// SinkConfig describes a single output sink.
type SinkConfig struct {
	Type   string `yaml:"type"`
	Target string `yaml:"target,omitempty"`
}

// FilterConfig describes a filter rule applied to the pipeline.
type FilterConfig struct {
	Type    string `yaml:"type"`
	Value   string `yaml:"value"`
	Level   string `yaml:"level,omitempty"`
}

// Config is the top-level logpipe configuration.
type Config struct {
	Sources []SourceConfig `yaml:"sources"`
	Sinks   []SinkConfig   `yaml:"sinks"`
	Filters []FilterConfig `yaml:"filters"`
	Workers int            `yaml:"workers"`
	Metrics struct {
		Enabled  bool   `yaml:"enabled"`
		Interval int    `yaml:"interval_seconds"`
	} `yaml:"metrics"`
	Health struct {
		Addr string `yaml:"addr"`
	} `yaml:"health"`
}

// Load reads and validates a Config from the YAML file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse: %w", err)
	}
	if len(cfg.Sources) == 0 {
		return nil, errors.New("config: at least one source is required")
	}
	if len(cfg.Sinks) == 0 {
		return nil, errors.New("config: at least one sink is required")
	}
	for i, s := range cfg.Sources {
		if s.Name == "" {
			return nil, fmt.Errorf("config: source[%d]: name is required", i)
		}
		if s.Path == "" {
			return nil, fmt.Errorf("config: source[%d]: path is required", i)
		}
	}
	if cfg.Workers <= 0 {
		cfg.Workers = defaultWorkers
	}
	return &cfg, nil
}
