// Package http is the package to get status though HTTP.
package http

import (
	"context"
	"fmt"
	"io"
	"log"
	nethttp "net/http"
	neturl "net/url"
	"strconv"
	"time"

	"github.com/rangzen/otel-status/package/status"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const name = "http"

// HTTP is the main structure to use HTTP status.
type HTTP struct {
	SC     status.Config
	Method string
	Site   string
}

// Config returns the status.Config of the HTTP status.
func (h HTTP) Config() status.Config {
	return h.SC
}

// State do the traces about the HTTP status.
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/http.md
func (h HTTP) State(tracer trace.Tracer) error {
	start := time.Now()
	_, span := tracer.Start(context.Background(), "GET "+h.Site,
		trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	// Defers use LIFO, so the defer that ends the span should be called first.
	defer func(Start time.Time) {
		span.SetAttributes(
			attribute.Int64("duration", time.Since(start).Milliseconds()),
		)
	}(start)

	url, err := neturl.Parse(h.Site)
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

	log.Printf("%s: %q with status %q", name, url.Host, strconv.Itoa(res.StatusCode))
	span.SetAttributes(
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
