# Documentation

## Example of usage

### Uptrace

![Uptrace](Screenshot_Uptrace.png)

## Examples

### Docker Compose with OpenTelemetry Collector and Uptrace

See the [compose-otelcol-uptrace](./examples/compose-otelcol-uptrace) directory for files.

This example use [otel-status](https://github.com/open-telemetry/opentelemetry-collector) and the [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector) to send spans and metrics to [Uptrace](https://uptrace.dev).

See https://stackoverflow.com/a/48547074 for the `172.17.0.1` trick.