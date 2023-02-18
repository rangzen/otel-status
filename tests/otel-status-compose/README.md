# Local Docker Compose setup for testing

This directory contains a Docker Compose configuration for running the
tests locally. It is not used in CI.

## Requirements

- Docker
- Docker Compose

## Services

The following services are started:

- `service-http`: local HTTP server, configured to respond to requests with
  with different status codes depending on the port used. 
- `otel-collector`: OpenTelemetry Collector configured 
  to receive and export traces and metrics.
- `clickhouse`: ClickHouse database, configured as backend for Uptrace.
- `uptrace`: Uptrace server, configured to receive traces and metrics from
  the OpenTelemetry Collector and directly.

## Usage

To run the services, run the following command:
```shell
docker-otel-status-compose up -d
```
