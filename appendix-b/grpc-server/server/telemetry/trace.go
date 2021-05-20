package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"go.opentelemetry.io/otel/exporters/jaeger"
)

func InitTracing(jaegerAddr string) error {

	/*traceExporter, err := stdouttrace.New()
	if err != nil {
		return tr, nil, err
	}*/

	traceExporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(jaegerAddr + "/api/traces"),
		),
	)
	if err != nil {
		return err
	}

	// For production uncomment next line, and comment the liner after, after
	// bsp := sdktrace.NewBatchSpanProcessor(traceExporter)

	// only recommended for demos/debugging
	bsp := sdktrace.NewSimpleSpanProcessor(traceExporter)

	// Semantic conventions:
	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/README.md
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("PkgServer"),
			),
		),
	)
	otel.SetTracerProvider(tp)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
	)
	otel.SetTextMapPropagator(propagator)
	return nil
}
