package metric

import (
	"context"
	"math/rand"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const metricName = "provider_status_manager_of_provider_request_total"

func Count() {
	counter, err := otel.GetMeterProvider().
		Meter(meterName).
		Int64Counter(metricName, metric.WithDescription("Provider Requests"))

	if err != nil {
		panic(err)
	}

	attr := []attribute.KeyValue{
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

func randomProvider() int {
	random := rand.Intn(5)
	return random + 1000
}

func randomStatusCode() int {
	codes := []int{200, 400, 500}
	return codes[rand.Intn(len(codes))]
}
