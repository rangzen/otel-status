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
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

const pluginName = "http"

const (
	otelStatusHTTPDuration = "otel.status.http.duration"
	otelStatusHTTPName     = "otel.status.http.name"
)

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
	_, span := tracer.Start(ctx, "GET "+h.URL,
		trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	// Defers use LIFO, so the defer that ends the span should be called first.
	defer func(Start time.Time) {
		// Span
		elapsedTime := time.Since(start).Milliseconds()
		span.SetAttributes(
			attribute.Int64("duration", elapsedTime),
		)

		// Metric
		opHistogram, err := meter.Int64Histogram(
			otelStatusHTTPDuration,
			instrument.WithUnit("ms"),
			instrument.WithDescription("Duration of the HTTP request"),
		)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		opHistogram.Record(ctx, elapsedTime,
			attribute.String(otelStatusHTTPName, h.SC.Name),
		)
	}(start)

	url, err := neturl.Parse(h.URL)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("parsing url: %w", err)
	}

	req, err := nethttp.NewRequest(h.Method, url.String(), nethttp.NoBody)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("creating HTTP request: %w", err)
	}
	client := &nethttp.Client{}
	res, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("doing HTTP %s: %w", h.Method, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}(res.Body)

	slog.Info("status", slog.String("plugin", pluginName), slog.String("url", url.String()), slog.String("status", strconv.Itoa(res.StatusCode)))
	span.SetAttributes(
		attribute.String(status.OtelStatusPluginName, pluginName),
		semconv.HTTPMethodKey.String("GET"),
		semconv.HTTPStatusCodeKey.Int(res.StatusCode),
		semconv.HTTPFlavorKey.String("1.1"),
		semconv.HTTPURLKey.String(url.String()),
		semconv.NetPeerNameKey.String(url.Host),
		semconv.NetPeerPortKey.String(url.Port()),
	)
	if res.StatusCode >= 400 {
		span.SetStatus(codes.Error, fmt.Sprintf("HTTP status code %d", res.StatusCode))
	}

	return nil
}
