// main otel-status is the main command and the main example of how to use this library.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
)

var tracer = otel.Tracer("github.com/rangzen/otel-status")
var meter = global.MeterProvider().Meter("github.com/rangzen/otel-status")

func main() {
	// Load configuration file.
	configPath := flag.String("config", "", "Path to the configuration file.")
	flag.Parse()

	if *configPath == "" {
		log.Fatal("You must provide a configuration file.")
	}

	conf, err := config.FromFile(*configPath)
	if err != nil {
		log.Fatal(fmt.Errorf("loading configuration file: %w", err))
	}

	// Prepare connexion to Open Telemetry Traces.
	if err = initTracer(); err != nil {
		log.Fatal(err)
	}

	// Prepare connexion to Open Telemetry Metrics.
	if err = initMeter(); err != nil {
		log.Fatal(err)
	}

	// Print all OTEL_ environment variables.
	log.Println("Open Telemetry environment variables:")
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "OTEL_") {
			log.Println(e)
		}
	}

	// Cron all status on local time zone.
	scheduler := gocron.NewScheduler(time.Local)
	for _, s := range conf.States.HTTP {
		stater := http.HTTP{
			SC: status.Config{
				Name:        s.Name,
				Description: s.Description,
				Cron:        s.Cron,
			},
			Method: s.Method,
			URL:    s.URL,
		}
		log.Printf("Scheduling %q at %q", stater.Config().Name, stater.Config().Cron)
		if stater.Config().IsDuration() {
			_, err = scheduler.Every(stater.Config().CronDuration()).Do(stater.State, tracer, meter)
		} else {
			_, err = scheduler.Cron(stater.Config().CronExp()).Do(stater.State, tracer, meter)
		}
		if err != nil {
			log.Fatal(fmt.Errorf("scheduling %s: %w", stater.Config().Name, err))
		}
	}
	log.Println("Status scheduled:", scheduler.Len())
	scheduler.StartBlocking()
}

// initTracer prepares connexion to Open Telemetry Traces.
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

// initMeter prepares connexion to Open Telemetry Metrics.
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
