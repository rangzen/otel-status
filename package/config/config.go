// Package config provides the configuration for the application.
package config

import (
	"fmt"
	"os"

	"github.com/rangzen/otel-status/package/status/http"
	"gopkg.in/yaml.v3"
)

// Config is the configuration root type for configuration file.
type Config struct {
	States States `yaml:"states"`
}

// States is the configuration for all the status.
type States struct {
	HTTP []http.Config `yaml:"http"`
}

// FromBytes returns the States from the given slice of bytes.
func FromBytes(data []byte) (Config, error) {
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("unmarshaling config file: %w", err)
	}
	return config, nil
}

// FromFile returns the States from the given file path.
func FromFile(path string) (Config, error) {
	configData, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading config file: %w", err)
	}
	return FromBytes(configData)
}
