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

// main otel-status is the main command and the main example of how to use this library.
package main

import (
	"context"
	"flag"
	"fmt"
	neturl "net/url"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rangzen/otel-status/package/config"
	"github.com/rangzen/otel-status/package/status"
	"github.com/rangzen/otel-status/package/status/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"golang.org/x/exp/slog"
)

const (
	instrumentName = "github.com/rangzen/otel-status"
)

func main() {
	slog.Info("starting otel-status")

	// Load configuration file.
	configPath := flag.String("config", "", "Path to the configuration file.")
	flag.Parse()

	if *configPath == "" {
		slog.Error("You must provide a configuration file. Try -h for help.", nil)
		os.Exit(1)
	}

	slog.Info("loading configuration", "path", *configPath)
	conf, err := config.FromFile(*configPath)
	if err != nil {
		slog.Error("loading configuration", err, "path", *configPath)
		os.Exit(1)
	}

	// Prepare connection to Open Telemetry Traces.
	if err = initTracer(); err != nil {
		slog.Error("initializing tracer", err)
		os.Exit(1)
	}

	// Prepare connection to Open Telemetry Metrics.
	if err = initMeter(); err != nil {
		slog.Error("initializing meter", err)
		os.Exit(1)
	}

	// Print all OTEL_ environment variables.
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "OTEL_") {
			split := strings.SplitN(e, "=", 2)
			slog.Info("Open Telemetry environment variable", "name", split[0], "value", split[1])
		}
	}

	// Cron all status on local time zone.
	var tracer = otel.Tracer(instrumentName)
	var meter = global.MeterProvider().Meter(instrumentName)
	scheduler := gocron.NewScheduler(time.Local)
	for _, s := range conf.States.HTTP {
		// Parse the URL.
		url, err := neturl.Parse(s.URL)
		if err != nil {
			slog.Error("parsing URL", err, "url", s.URL)
			continue
		}

		stater := http.HTTP{
			SC: status.Config{
				Name:        s.Name,
				Description: s.Description,
				Cron:        s.Cron,
			},
			Method: s.Method,
			URL:    url,
			Values: s.Values,
		}
		slog.Info("scheduling", "plugin", http.PluginName, "name", stater.Config().Name, "cron", stater.Config().Cron)
		if stater.Config().IsDuration() {
			_, err = scheduler.Every(stater.Config().CronDuration()).Do(stater.State, tracer, meter)
		} else {
			_, err = scheduler.Cron(stater.Config().CronExp()).Do(stater.State, tracer, meter)
		}
		if err != nil {
			slog.Error("scheduling", err, "plugin", http.PluginName, "name", stater.Config().Name)
			os.Exit(1)
		}
	}
	slog.Info("scheduled", "count", scheduler.Len())
	scheduler.StartBlocking()
}

// initTracer prepares connection to Open Telemetry Traces.
// All the configuration is done via environment variables.
func initTracer() error {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(),
	)
	if err != nil {
		return fmt.Errorf("creating Open Telemetry traces exporter: %w", err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("otel-status"),
		),
	)
	if err != nil {
		return fmt.Errorf("creating Open Telemetry traces resources: %w", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
		trace.WithResource(resources),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return nil
}

// initMeter prepares connection to Open Telemetry Metrics.
// All the configuration is done via environment variables.
func initMeter() error {
	exporter, err := otlpmetricgrpc.New(
		context.Background(),
	)
	if err != nil {
		return fmt.Errorf("creating Open Telemetry metrics exporter: %w", err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("otel-status"),
		),
	)
	if err != nil {
		return fmt.Errorf("creating Open Telemetry metrics resources: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(resources),
	)

	global.SetMeterProvider(meterProvider)

	return nil
}
