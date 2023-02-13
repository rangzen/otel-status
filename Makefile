.PHONY: lint
lint:
	revive -config revive.toml ./...

.PHONY: test-force-rebuild-service-http
test-force-rebuild-service-http:
	cd tests/compose && docker compose build --no-cache service-http

.PHONY: test-run-compose
test-run-compose:
	cd tests/compose && docker compose up -d

.PHONY: test-run
test-run:
	export OTEL_EXPORTER_OTLP_ENDPOINT=grpc://localhost:14317 && \
	export OTEL_EXPORTER_OTLP_INSECURE=true && \
	export OTEL_EXPORTER_OTLP_TRACES_HEADERS=UPTRACE-DSN=http://project2_secret_token@localhost:14317/2 && \
	export OTEL_RESOURCE_ATTRIBUTES=deployment.environment=dev && \
	go run tests/compose/cmd/otel-status.go
