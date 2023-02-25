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

// Package status regroups everything related to the different status that can be checked.
package status

import (
	"strings"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	// OtelStatusPluginName is the key for the plugin name.
	OtelStatusPluginName = "otelstatus.plugin.name"
)

// Stater is the interface that wraps the Config methods.
type Stater interface {
	Config() Config
	State(tracer trace.Tracer, meter metric.Meter) error
}

// Config is the main structure to use status.
type Config struct {
	Name        string
	Description string
	Cron        string
}

// CronExp returns the cron expression.
func (s Config) CronExp() string {
	return s.Cron
}

// IsDuration returns true if the cron is of type time.Duration.
// The cron is of type time.Duration if it starts with @.
func (s Config) IsDuration() bool {
	return strings.HasPrefix(s.Cron, "@")
}

// CronDuration returns the time.Duration of the cron.
func (s Config) CronDuration() string {
	return s.Cron[1:]
}
