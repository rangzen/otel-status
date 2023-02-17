// Package http is the package to get status though HTTP.
package http

import (
	"context"
	"fmt"
	"io"
	nethttp "net/http"
	neturl "net/url"
	"strconv"
	"time"

	"github.com/rangzen/otel-status/package/status"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

// PluginName is the name of the plugin.
const PluginName = "http"

const (
	otelStatusHTTPName     = "otel.status.http.name"
	otelStatusHTTPDuration = "otel.status.http.duration"
	otelStatusHTTPStatus   = "otel.status.http.status"
)

var httpStatusClass = [5]string{"1xx", "2xx", "3xx", "4xx", "5xx"}

// Config is the configuration for an HTTP status.
type Config struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Cron        string `yaml:"cron" default:"@10m"`
	Method      string `yaml:"method" default:"GET"`
	URL         string `yaml:"url"`
}

// HTTP is the main structure to use HTTP status.
type HTTP struct {
	SC     status.Config
	Method string
	URL    string
}

// Config returns the status.Config of the HTTP status.
func (h HTTP) Config() status.Config {
	return h.SC
}

// State do the traces about the HTTP status.
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/http.md
func (h HTTP) State(tracer trace.Tracer, meter metric.Meter) error {
	ctx := context.Background()
	start := time.Now()

	// Create a span.
	_, span := tracer.Start(ctx, fmt.Sprintf("%s %s", h.Method, h.URL),
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String(status.OtelStatusPluginName, PluginName),
			semconv.HTTPMethodKey.String(h.Method),
		),
	)
	// Defers use LIFO, so the defer that ends the span
	// should be the first defer to then be called last.
	defer span.End()

	// Parse the URL.
	url, err := neturl.Parse(h.URL)
	if err != nil {
		return errorHandling(err, "parsing url", span)
	}
	span.SetAttributes(
		semconv.HTTPSchemeKey.String(url.Scheme),
		semconv.HTTPURLKey.String(url.String()),
		// As server
		semconv.NetHostNameKey.String(url.Hostname()),
		semconv.NetHostPortKey.String(url.Port()),
		// As client
		semconv.NetPeerNameKey.String(url.Hostname()),
		semconv.NetPeerPortKey.String(url.Port()),
	)

	// Create the HTTP request.
	req, err := nethttp.NewRequest(h.Method, url.String(), nethttp.NoBody)
	if err != nil {
		return errorHandling(err, "creating HTTP request", span)
	}

	// Do the HTTP request.
	client := &nethttp.Client{}
	res, err := client.Do(req)
	if err != nil {
		return errorHandling(err, "doing HTTP client", span)
	}
	defer func(body io.ReadCloser) {
		// Close the response body.
		err = body.Close()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}(res.Body)

	slog.Info("status", slog.String("plugin", PluginName), slog.String("url", url.String()), slog.String("status", strconv.Itoa(res.StatusCode)))

	// Record the status.
	span.SetAttributes(
		semconv.HTTPStatusCodeKey.Int(res.StatusCode),
	)
	if res.StatusCode >= 400 {
		span.SetStatus(codes.Error, fmt.Sprintf("HTTP status code %d", res.StatusCode))
	}

	// Record the duration.
	elapsedTime := time.Since(start).Milliseconds()
	span.SetAttributes(
		attribute.Int64("duration", elapsedTime),
	)

	// Record the duration in the meter.
	durationMetric, err := meter.Int64Histogram(
		otelStatusHTTPDuration,
		instrument.WithUnit(unit.Milliseconds),
		instrument.WithDescription("Duration of the HTTP request"),
	)
	if err != nil {
		return errorHandling(err, "creating HTTP request duration metric", span)
	}
	durationMetric.Record(ctx, elapsedTime, attribute.String(otelStatusHTTPName, h.SC.Name))

	// Record the family status as a compromise between the number of metrics and the number of labels in the meter.
	_, err = meter.Int64ObservableGauge(
		otelStatusHTTPStatus,
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Status of the HTTP request"),
		instrument.WithInt64Callback(func(_ context.Context, observer instrument.Int64Observer) error {
			statusClassIndex := (res.StatusCode / 100) - 1
			for i := 0; i < len(httpStatusClass); i++ {
				var value int64
				if i == statusClassIndex {
					value = 1
				}
				observer.Observe(
					value,
					attribute.String(otelStatusHTTPName, h.SC.Name),
					semconv.HTTPURLKey.String(url.String()),
					semconv.HTTPStatusCodeKey.Int(res.StatusCode),
					semconv.HTTPMethodKey.String(h.Method),
					attribute.String("http.status_class", httpStatusClass[i]),
				)
			}
			return nil
		}),
	)
	if err != nil {
		return errorHandling(err, "creating HTTP request status metric", span)
	}

	return nil
}

// errorHandling is a helper function to handle errors.
// It logs the error, records it in the span and returns it.
func errorHandling(err error, msg string, span trace.Span) error {
	e := fmt.Errorf("%s: %w", msg, err)
	slog.Error(msg, e, slog.String("plugin", PluginName))
	span.RecordError(e)
	span.SetStatus(codes.Error, e.Error())
	return e
}
