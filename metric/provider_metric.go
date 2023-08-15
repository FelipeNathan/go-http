package metric

import (
	"context"
	"fmt"
	"math/rand"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	metricName = "provider_status_manager_of_provider_request_total"
	gaugeName  = "provider_status_manager_of_provider_current_status"
)

var status = []string{"DISABLED", "TEMPORARILY_UNAVAILABLE", "ENABLED"}

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
	index, _ := randomStatus()
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
				fmt.Printf("%v - %v\n", index, provider)
				io.Observe(index, metric.WithAttributes(attr...))
				return nil
			}),
		)

	if err != nil {
		panic(err)
	}
}

func randomProvider() int {
	random := rand.Intn(5)
	return random + 1000
}

func randomStatusCode() int {
	codes := []int{200, 400, 500}
	return codes[rand.Intn(len(codes))]
}

func randomStatus() (int64, string) {
	index := rand.Intn(len(status))
	return int64(index), status[index]
}
