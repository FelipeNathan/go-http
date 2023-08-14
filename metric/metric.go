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

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("my_service"),
		semconv.ServiceVersionKey.String("v0.0.0"),
	)

	otlpExporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint("localhost:4318"),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}

	stdOutExporter, err := stdoutmetric.New()
	if err != nil {
		panic(err)
	}

	fiveSeconds := metricSdk.WithInterval(time.Second * 5)
	mp := metricSdk.NewMeterProvider(
		metricSdk.WithResource(resources),
		metricSdk.WithReader(metricSdk.NewPeriodicReader(otlpExporter, fiveSeconds)),
		metricSdk.WithReader(metricSdk.NewPeriodicReader(stdOutExporter, fiveSeconds)),
	)

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
