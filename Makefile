.PHONY: lint
lint:
	revive -config revive.toml ./...

.PHONY: test-force-rebuild-service-http
test-force-rebuild-service-http:
	cd tests/otel-status-compose && docker compose build --no-cache service-http

.PHONY: test-compose-up
test-compose-up:
	cd tests/otel-status-compose && docker compose up -d

.PHONY: test-run-cmd-collector
test-run-cmd-collector:
	export OTEL_EXPORTER_OTLP_ENDPOINT=grpc://localhost:4317 && \
	export OTEL_EXPORTER_OTLP_INSECURE=true && \
	export OTEL_RESOURCE_ATTRIBUTES=deployment.environment=dev && \
	go run cmd/otel-status/otel-status.go -config tests/otel-status-compose/otel-status.yaml

.PHONY: test-run-cmd-uptrace
test-run-cmd:
	export OTEL_EXPORTER_OTLP_ENDPOINT=grpc://localhost:14317 && \
	export OTEL_EXPORTER_OTLP_METRICS_ENDPOINT=grpc://localhost:14317 && \
	export OTEL_EXPORTER_OTLP_INSECURE=true && \
	export OTEL_EXPORTER_OTLP_TRACES_HEADERS=UPTRACE-DSN=http://project2_secret_token@localhost:14317/2 && \
	export OTEL_EXPORTER_OTLP_METRICS_HEADERS=UPTRACE-DSN=http://project2_secret_token@localhost:14317/2 && \
	export OTEL_RESOURCE_ATTRIBUTES=deployment.environment=dev && \
	go run cmd/otel-status/otel-status.go -config tests/otel-status-compose/otel-status.yaml

.PHONY: generate-changelog-dry
generate-changelog-dry:
	conventional-changelog -p conventionalcommits -i CHANGELOG.md
