/*******************************************************************************
 * Copyright (c) 2023 Cedric L'homme.
 *
 * This file is part of otel-status.
 *
 * otel-status is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or (at your option) any later version.
 *
 *  otel-status is distributed in the hope that it will be useful, but
 *  WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *  See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with otel-status. If not, see <https://www.gnu.org/licenses/>.
 ******************************************************************************/

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
