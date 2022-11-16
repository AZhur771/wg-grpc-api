package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseConfig(path string) (*Config, error) {
	config := NewConfig()

	f, err := os.Open(path)
	if err != nil {
		return config, fmt.Errorf("error while openening config file %s: %w", path, err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)

	if err := yaml.Unmarshal(b, config); err != nil {
		return config, fmt.Errorf("error while unmarshalling config file %s: %w", path, err)
	}

	return config, nil
}
