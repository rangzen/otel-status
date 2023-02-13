// Package status regroups everything related to the different status that can be checked.
package status

import (
	"strings"

	"go.opentelemetry.io/otel/trace"
)

// Stater is the interface that wraps the Config methods.
type Stater interface {
	Config() Config
	State(tracer trace.Tracer) error
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
func (s Config) IsDuration() bool {
	return strings.HasPrefix(s.Cron, "@")
}

// CronDuration returns the time.Duration of the cron.
func (s Config) CronDuration() string {
	return s.Cron[1:]
}