// main otel-status is the main command and the main example of how to use this library.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rangzen/otel-status/package/status"
	statushttp "github.com/rangzen/otel-status/package/status/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
)

var tracer = otel.Tracer("github.com/rangzen/otel-status")

func main() {
	// Prepare connexion to Open Telemetry Traces.
	if err := initTracer(); err != nil {
		log.Fatal(err)
	}

	// Print all OTEL_ environment variables.
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "OTEL_") {
			log.Printf("%s\n", e)
		}
	}

	// TODO Get these from configuration.
	var staters = []status.Stater{
		statushttp.HTTP{
			SC: status.Config{
				Name:        "HTTP GET on localhost:8080 (HTML)",
				Description: "Check if the HTTP GET is working.",
				Cron:        "@4s",
			},
			Method: http.MethodGet,
			Site:   "http://localhost:8080",
		},
		statushttp.HTTP{
			SC: status.Config{
				Name:        "HTTP HEAD on localhost:8081 (slow, JSON)",
				Description: "Check if the HTTP HEAD is working.",
				Cron:        "@7s",
			},
			Method: http.MethodHead,
			Site:   "http://localhost:8081",
		},
		statushttp.HTTP{
			SC: status.Config{
				Name:        "HTTP HEAD on localhost:8082 (401)",
				Description: "Check if the HTTP HEAD is working.",
				Cron:        "*/2 * * * *",
			},
			Method: http.MethodHead,
			Site:   "http://localhost:8082",
		},
	}

	// Cron all status.
	scheduler := gocron.NewScheduler(time.UTC)
	for _, stater := range staters {
		var err error
		if stater.Config().IsDuration() {
			_, err = scheduler.Every(stater.Config().CronDuration()).Do(stater.State, tracer)
		} else {
			_, err = scheduler.Cron(stater.Config().CronExp()).Do(stater.State, tracer)
		}
		if err != nil {
			log.Fatal(fmt.Errorf("scheduling %s: %w", stater.Config().Name, err))
		}
	}
	log.Println(scheduler.Len(), "status scheduled.")
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
