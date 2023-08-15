package metric

import (
	"context"
	"math/rand"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	metricName    = "provider_status_manager_of_provider_request_total"
	gaugeName     = "provider_status_manager_of_provider_current_latency_request"
	histogramName = "provider_status_manager_of_provider_latency_request"
)

type kv []attribute.KeyValue

func Count() {
	counter, err := otel.GetMeterProvider().
		Meter(meterName).
		Int64Counter(metricName, metric.WithDescription("Provider Requests"))

	if err != nil {
		panic(err)
	}

	attr := kv{
		attribute.Int("providerId", randomProvider()),
		attribute.Int("httpStatusCode", randomStatusCode()),
		attribute.String("resourceConsumptionContext", "AISP_CONSENT"),
	}

	counter.Add(
		context.Background(),
		1,
		metric.WithAttributes(attr...),
	)
}

func Gauge() {
	latency := randomLatency()
	provider := randomProvider()

	attr := kv{
		attribute.Int("providerId", provider),
	}

	_, err := otel.GetMeterProvider().
		Meter(meterName).
		Int64ObservableGauge(
			gaugeName,
			metric.WithDescription("Status of Providers"),
			metric.WithInt64Callback(func(ctx context.Context, io metric.Int64Observer) error {
				io.Observe(latency, metric.WithAttributes(attr...))
				return nil
			}),
		)

	if err != nil {
		panic(err)
	}
}

func Histogram() {
	histo, err := otel.GetMeterProvider().
		Meter(meterName).
		Int64Histogram(histogramName, metric.WithDescription("Latency of the providers requests"))

	if err != nil {
		panic(err)
	}

	attr := kv{
		attribute.Int("providerId", randomProvider()),
	}

	histo.Record(context.Background(), randomLatency(), metric.WithAttributes(attr...))
}

func randomProvider() int {
	random := rand.Intn(5)
	return random + 1000
}

func randomStatusCode() int {
	codes := []int{200, 400, 500}
	return codes[rand.Intn(len(codes))]
}

func randomLatency() int64 {
	ls := []int64{30, 60, 500, 700, 1200}
	return ls[rand.Intn(len(ls))]
}
