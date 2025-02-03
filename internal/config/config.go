package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	configPath = "config.yml"
)

func parseFile(path string, dest any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.New("failed to read config.")
	}

	if err := yaml.Unmarshal(data, dest); err != nil {
		return errors.New("failed to unmarshal config into struct.")
	}

	return nil
}
