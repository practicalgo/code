module github.com/practicalgo/code/appendix-a/http-server

go 1.16

require (
	github.com/DataDog/datadog-go v4.8.1+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/practicalgo/code/appendix-a/grpc-server/service v0.0.0
	github.com/rs/zerolog v1.24.0
	github.com/testcontainers/testcontainers-go v0.11.1
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.22.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.23.0
	go.opentelemetry.io/otel v1.0.0-RC3
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.0-RC3
	go.opentelemetry.io/otel/sdk v1.0.0-RC3
	go.opentelemetry.io/otel/trace v1.0.0-RC3
	gocloud.dev v0.23.0
	google.golang.org/grpc v1.39.0
)

replace github.com/practicalgo/code/appendix-a/grpc-server/service => ../grpc-server/service

// Remove replace and upgrade library once
// https://github.com/testcontainers/testcontainers-go/pull/342 is merged
// The tag used here is on my personal fork containing the change in PR:
// https://github.com/amitsaha/testcontainers-go/releases/tag/v0.11.1-pr-342
replace github.com/testcontainers/testcontainers-go => github.com/amitsaha/testcontainers-go v0.11.1-pr-342
