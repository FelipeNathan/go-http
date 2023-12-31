package jaeger

import (
	"context"
	"fmt"

	"github.com/FelipeNathan/go-http/internal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func Config() {

	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(internal.AppName),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(r),
		sdktrace.WithBatcher(newExporter()),
	)

	otel.SetTracerProvider(tp)
}

func Shutdown() {
	fmt.Print("Shutting down trace provider")
	err := otel.GetTracerProvider().(*sdktrace.TracerProvider).Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

func newExporter() *otlptrace.Exporter {

	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint("localhost:4320"),
		otlptracehttp.WithInsecure(),
	)

	if exporter, err := otlptrace.New(context.Background(), client); err != nil {
		panic(err)
	} else {
		return exporter
	}
}

func Trace(ctx context.Context, spanName string, attr ...attribute.KeyValue) (context.Context, trace.Span) {
	return otel.Tracer(internal.AppName).Start(ctx, spanName, trace.WithAttributes(attr...))
}
