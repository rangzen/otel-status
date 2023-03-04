# Build this Dockerfile from the root of the repository.
FROM golang AS builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /go/bin/otel-status ./cmd/otel-status

# TODO: Use distroless/static:nonroot once it's available.
#FROM grc.io/distroless/static-debian11 AS runtime
FROM golang AS runtime

COPY --from=builder /go/bin/otel-status /

# USER nonroot:nonroot

CMD [ "/otel-status", "-config", "/otel-status.yaml" ]
