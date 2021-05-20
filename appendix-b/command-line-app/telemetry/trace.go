package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/exporters/jaeger"
)

type TraceReporter struct {
	Client trace.Tracer
	Ctx    context.Context
}

func InitTracing(jaegerAddr, version string) (TraceReporter, *sdktrace.TracerProvider, error) {

	var tr TraceReporter

	/*traceExporter, err := stdouttrace.New()
	if err != nil {
		return tr, nil, err
	}*/

	traceExporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(jaegerAddr),
		),
	)
	if err != nil {
		return tr, nil, err
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
				semconv.ServiceNameKey.String("PkgServer-Cli"),
			),
		),
	)
	otel.SetTracerProvider(tp)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
	)
	otel.SetTextMapPropagator(propagator)

	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/baggage/api.md
	v1, err := baggage.NewMember("cli_version", version)
	if err != nil {
		return tr, tp, err
	}
	bag, err := baggage.New(v1)
	if err != nil {
		return tr, tp, err
	}

	tr.Client = otel.Tracer("")
	ctx := context.Background()
	tr.Ctx = baggage.ContextWithBaggage(ctx, bag)

	return tr, tp, nil
}
