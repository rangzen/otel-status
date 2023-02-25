package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/rangzen/otel-status/package/status"
	otelhttp "github.com/rangzen/otel-status/package/status/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestHTTP_Status(t *testing.T) {
	t.Run("a 200 status, should create a span without an error status", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		urlParsed, err := url.Parse(mockServer.URL)
		require.NoError(t, err)

		exp := tracetest.NewInMemoryExporter()

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
		)

		mockTracer := tp.Tracer("test-tracer")

		rdr := metric.NewManualReader()
		mp := metric.NewMeterProvider(metric.WithReader(rdr))
		mockMeter := mp.Meter("test-meter")

		stater := otelhttp.HTTP{
			SC: status.Config{
				Name:        "Test",
				Description: "Test 200",
				Cron:        "@99m",
			},
			Method: http.MethodGet,
			URL:    urlParsed,
			Values: nil,
		}

		err = stater.State(mockTracer, mockMeter)
		require.NoError(t, err)

		ctx := context.Background()
		// Assert span
		err = tp.ForceFlush(ctx)
		require.NoError(t, err)

		spans := exp.GetSpans()
		require.Len(t, spans, 1)
		require.Equal(t, codes.Unset, spans[0].Status.Code)

		// Assert metric
		m, err := rdr.Collect(ctx)
		assert.NoError(t, err)

		require.Len(t, m.ScopeMetrics, 1)
		require.Len(t, m.ScopeMetrics[0].Metrics, 2)
	})

	t.Run("a 401 error, should create a span with an error status", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer mockServer.Close()

		urlParsed, err := url.Parse(mockServer.URL)
		require.NoError(t, err)

		exp := tracetest.NewInMemoryExporter()

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
		)

		mockTracer := tp.Tracer("test-tracer")

		rdr := metric.NewManualReader()
		mp := metric.NewMeterProvider(metric.WithReader(rdr))
		mockMeter := mp.Meter("test-meter")

		stater := otelhttp.HTTP{
			SC: status.Config{
				Name:        "Test",
				Description: "Test 401",
				Cron:        "@99m",
			},
			Method: http.MethodGet,
			URL:    urlParsed,
			Values: nil,
		}

		err = stater.State(mockTracer, mockMeter)
		require.NoError(t, err)

		// Assert span
		ctx := context.Background()
		err = tp.ForceFlush(ctx)
		require.NoError(t, err)

		spans := exp.GetSpans()
		require.Len(t, spans, 1)
		require.Equal(t, codes.Error, spans[0].Status.Code)

		// Assert metric
		m, err := rdr.Collect(ctx)
		assert.NoError(t, err)

		require.Len(t, m.ScopeMetrics, 1)
		require.Len(t, m.ScopeMetrics[0].Metrics, 2)
	})

	t.Run("an HTTP error, should create a span with an error status", func(t *testing.T) {
		urlParsed, err := url.Parse("https://tralalala-tsouintsouin-les-ptites-boules---aliiiaaaaaa-prt.dev")
		require.NoError(t, err)

		exp := tracetest.NewInMemoryExporter()
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
		)
		mockTracer := tp.Tracer("test-tracer")

		rdr := metric.NewManualReader()
		mp := metric.NewMeterProvider(metric.WithReader(rdr))
		mockMeter := mp.Meter("test-meter")

		stater := otelhttp.HTTP{
			SC: status.Config{
				Name:        "Test",
				Description: "Test non-existent domain",
				Cron:        "@99m",
			},
			Method: http.MethodGet,
			URL:    urlParsed,
			Values: nil,
		}

		err = stater.State(mockTracer, mockMeter)
		require.Error(t, err)

		ctx := context.Background()
		// Assert span
		err = tp.ForceFlush(ctx)
		require.NoError(t, err)

		spans := exp.GetSpans()
		require.Len(t, spans, 1)
		require.Equal(t, codes.Error, spans[0].Status.Code)

		// Assert metric
		m, err := rdr.Collect(ctx)
		assert.NoError(t, err)

		require.Len(t, m.ScopeMetrics, 1)
		require.Len(t, m.ScopeMetrics[0].Metrics, 1)
	})
}
