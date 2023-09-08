package metric

import (
	"context"
	"fmt"
	"time"

	"github.com/FelipeNathan/go-http/internal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	metricSdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// https://opentelemetry.io/docs/instrumentation/go/exporters/
// https://grafana.com/docs/opentelemetry/instrumentation/go/manual-instrumentation/
const (
	meterName = "github.com/FelipeNathan/go-http"
)

type meterOptions []metricSdk.Option

func Config() {

	options := meterOptions{}.
		withResource().
		withReaders()

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

func (options meterOptions) withResource() meterOptions {
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(internal.AppName),
		semconv.ServiceVersionKey.String("v0.0.0"),
	)

	return append(options, metricSdk.WithResource(resources))
}

func (options meterOptions) withReaders() meterOptions {
	fiveSeconds := metricSdk.WithInterval(time.Second * 5)

	if otlpExporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint("localhost:4318"),
		otlpmetrichttp.WithInsecure(),
	); err != nil {
		panic(err)
	} else {
		r := metricSdk.WithReader(metricSdk.NewPeriodicReader(otlpExporter, fiveSeconds))
		options = append(options, r)
	}

	// if stdOutExporter, err := stdoutmetric.New(); err != nil {
	// 	panic(err)
	// } else {
	// 	r := metricSdk.WithReader(metricSdk.NewPeriodicReader(stdOutExporter, fiveSeconds))
	// 	options = append(options, r)
	// }

	return options
}
