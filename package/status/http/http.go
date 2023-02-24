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
	otelStatusHTTPName     = "otelstatus.http.name"
	otelStatusHTTPDuration = "otelstatus.http.duration"
	otelStatusHTTPError    = "otelstatus.http.error"
	otelStatusHTTPStatus   = "otelstatus.http.status"
)

var httpStatusClass = [5]string{"1xx", "2xx", "3xx", "4xx", "5xx"}

// Config is the configuration for an HTTP status.
type Config struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Cron        string `yaml:"cron" default:"@10m"`
	Method      string `yaml:"method" default:"GET"`
	URL         string `yaml:"url"`
	// Values is a map of key/value to add to the spans.
	Values map[string]string `yaml:"values"`
}

// HTTP is the main structure to use HTTP status.
type HTTP struct {
	SC     status.Config
	Method string
	URL    *neturl.URL
	Values map[string]string
	// previousClass is the previous state of the HTTP status class metric.
	previousClass [5]bool
}

// Config returns the status.Config of the HTTP status.
func (h *HTTP) Config() status.Config {
	return h.SC
}

// State do the traces about the HTTP status.
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/http.md
func (h *HTTP) State(tracer trace.Tracer, meter metric.Meter) error {
	ctx := context.Background()
	start := time.Now()

	// Prepare additional attributes from config.
	var valuesAttributes []attribute.KeyValue
	for k, v := range h.Values {
		valuesAttributes = append(valuesAttributes, attribute.String(k, v))
	}

	// Create a span.
	_, span := tracer.Start(ctx, fmt.Sprintf("%s %s", h.Method, h.URL),
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String(status.OtelStatusPluginName, PluginName),
			semconv.HTTPMethodKey.String(h.Method),
		),
		trace.WithAttributes(valuesAttributes...),
	)
	// Defers use LIFO, so the defer that ends the span
	// should be the first defer to then be called last.
	defer span.End()

	span.SetAttributes(
		semconv.HTTPSchemeKey.String(h.URL.Scheme),
		semconv.HTTPURLKey.String(h.URL.String()),
		// As server
		semconv.NetHostNameKey.String(h.URL.Hostname()),
		semconv.NetHostPortKey.String(h.URL.Port()),
		// As client
		semconv.NetPeerNameKey.String(h.URL.Hostname()),
		semconv.NetPeerPortKey.String(h.URL.Port()),
	)

	// Create the HTTP request.
	req, err := nethttp.NewRequest(h.Method, h.URL.String(), nethttp.NoBody)
	if err != nil {
		return h.errorHandling(ctx, err, "creating HTTP request", span, meter)
	}

	// Do the HTTP request.
	client := &nethttp.Client{}
	res, err := client.Do(req)
	if err != nil {
		return h.errorHandling(ctx, err, "doing HTTP client", span, meter)
	}
	defer func(body io.ReadCloser) {
		// Close the response body.
		err = body.Close()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}(res.Body)

	slog.Info("status", slog.String("plugin", PluginName), slog.String("url", h.URL.String()), slog.String("status", strconv.Itoa(res.StatusCode)))

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
		return h.errorHandling(ctx, err, "creating HTTP request duration metric", span, meter)
	}
	durationMetric.Record(ctx, elapsedTime,
		attribute.String(otelStatusHTTPName, h.SC.Name),
		semconv.HTTPURLKey.String(h.URL.String()),
	)

	// Record the family status as a compromise between the number of metrics and the number of labels in the meter.
	// We need to keep an internal state to know if the status has changed,
	// and calculate the value to add to mimic the 0 or 1 values.
	// We cannot achieve the same behaviour as httpcheck with the current SDK,
	// there is no Int64Gauge.
	statusMetric, err := meter.Int64UpDownCounter(
		otelStatusHTTPStatus,
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Status of the HTTP request"),
	)
	if err != nil {
		return h.errorHandling(ctx, err, "creating HTTP request status metric", span, meter)
	}
	statusClassIndex := (res.StatusCode / 100) - 1
	for i := 0; i < len(httpStatusClass); i++ {
		val := int64(0)
		switch {
		case h.previousClass[i] && i != statusClassIndex:
			val = -1
		case !h.previousClass[i] && i == statusClassIndex:
			val = 1
		}
		statusMetric.Add(ctx, val,
			attribute.String(otelStatusHTTPName, h.SC.Name),
			semconv.HTTPURLKey.String(h.URL.String()),
			semconv.HTTPStatusCodeKey.Int(res.StatusCode),
			semconv.HTTPMethodKey.String(h.Method),
			attribute.String("http.status_class", httpStatusClass[i]))

		h.previousClass[i] = i == statusClassIndex
	}

	return nil
}

// errorHandling is a helper function to handle errors.
// It logs the error, records it in the span and returns it.
func (h *HTTP) errorHandling(ctx context.Context, err error, msg string, span trace.Span, meter metric.Meter) error {
	e := fmt.Errorf("%s: %w", msg, err)
	slog.Error(msg, e, slog.String("plugin", PluginName))
	span.RecordError(e)
	span.SetStatus(codes.Error, e.Error())

	// Record the metric error.
	errorMetric, err := meter.Int64Counter(
		otelStatusHTTPError,
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Error of the HTTP request"),
	)
	if err != nil {
		errorMetric.Add(ctx, 1,
			attribute.String(otelStatusHTTPName, h.SC.Name),
			semconv.HTTPURLKey.String(h.URL.String()),
			attribute.String("error.message", err.Error()),
		)
	}

	return e
}
