module github.com/practicalgo/code/appendix-a/grpc-server/server

go 1.16

require google.golang.org/grpc v1.39.0

require (
	github.com/DataDog/datadog-go v4.8.1+incompatible
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/practicalgo/code/appendix-a/grpc-server/service v0.0.0
	github.com/rs/zerolog v1.24.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.22.0
	go.opentelemetry.io/otel v1.0.0-RC3
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC3
	go.opentelemetry.io/otel/sdk v1.0.0-RC3
)

replace github.com/practicalgo/code/appendix-a/grpc-server/service => ../service
