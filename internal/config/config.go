package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SourceConfig defines a log source to tail.
type SourceConfig struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Tag  string `yaml:"tag"`
}

// SinkConfig defines a destination for log entries.
type SinkConfig struct {
	Name    string            `yaml:"name"`
	Type    string            `yaml:"type"` // stdout, file, http
	Options map[string]string `yaml:"options"`
}

// Config is the top-level logpipe configuration.
type Config struct {
	Sources []SourceConfig `yaml:"sources"`
	Sinks   []SinkConfig   `yaml:"sinks"`
	Workers int            `yaml:"workers"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	if cfg.Workers <= 0 {
		cfg.Workers = 4
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Sources) == 0 {
		return fmt.Errorf("at least one source must be defined")
	}
	if len(c.Sinks) == 0 {
		return fmt.Errorf("at least one sink must be defined")
	}
	for i, s := range c.Sources {
		if s.Path == "" {
			return fmt.Errorf("source[%d] missing path", i)
		}
		if s.Name == "" {
			return fmt.Errorf("source[%d] missing name", i)
		}
	}
	for i, s := range c.Sinks {
		if s.Type == "" {
			return fmt.Errorf("sink[%d] missing type", i)
		}
		if s.Name == "" {
			return fmt.Errorf("sink[%d] missing name", i)
		}
	}
	return nil
}
