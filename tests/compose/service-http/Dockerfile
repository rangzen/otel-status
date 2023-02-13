FROM golang:1.20 AS builder

WORKDIR /go/src/app

COPY ./go.mod .
COPY ./main.go .

RUN CGO_ENABLED=0 go build -o /go/bin/service-http main.go

# TODO: Use distroless/static:nonroot once it's available.
#FROM grc.io/distroless/static-debian11 AS runtime
FROM golang:1.20-alpine AS runtime

COPY --from=builder /go/bin/service-http /

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082

# USER nonroot:nonroot

ENTRYPOINT [ "/service-http" ]
