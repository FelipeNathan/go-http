package metric

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	metricSdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// https://opentelemetry.io/docs/instrumentation/go/exporters/
// https://grafana.com/docs/opentelemetry/instrumentation/go/manual-instrumentation/
const (
	meterName = "github.com/FelipeNathan/go-http"
)

func Config() {

	var options []metricSdk.Option
	options = withResource(options)
	options = withReaders(options)

	mp := metricSdk.NewMeterProvider(options...)

	otel.SetMeterProvider(mp)
}

func Shutdown() {
	fmt.Println("Shutting Down")
	mp := otel.GetMeterProvider()

	err := mp.(*metricSdk.MeterProvider).Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

func withResource(options []metricSdk.Option) []metricSdk.Option {
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("my_service"),
		semconv.ServiceVersionKey.String("v0.0.0"),
	)

	return append(options, metricSdk.WithResource(resources))
}

func withReaders(options []metricSdk.Option) []metricSdk.Option {
	fiveSeconds := metricSdk.WithInterval(time.Second * 5)

	if otlpExporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint("localhost:4318"),
		otlpmetrichttp.WithInsecure(),
	); err == nil {
		r := metricSdk.WithReader(metricSdk.NewPeriodicReader(otlpExporter, fiveSeconds))
		options = append(options, r)
	}

	if stdOutExporter, err := stdoutmetric.New(); err == nil {
		r := metricSdk.WithReader(metricSdk.NewPeriodicReader(stdOutExporter, fiveSeconds))
		options = append(options, r)
	}

	return options
}
