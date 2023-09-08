package metric

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	metricName    = "request_metric_request_total"
	gaugeName     = "request_metric_current_latency_request"
	histogramName = "request_metric_latency_request"
)

type kv []attribute.KeyValue

func Count(url string, statusCode int) {
	counter, err := otel.
		Meter(meterName).
		Int64Counter(metricName, metric.WithDescription("Url Requests"))

	if err != nil {
		panic(err)
	}

	attr := kv{
		attribute.String("url", url),
		attribute.Int("httpStatusCode", statusCode),
	}

	counter.Add(
		context.Background(),
		1,
		metric.WithAttributes(attr...),
	)
}

func Gauge(url string, latency int64) {

	attr := kv{
		attribute.String("url", url),
	}

	_, err := otel.
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

func Histogram(url string, latency int64) {
	histo, err := otel.
		Meter(meterName).
		Int64Histogram(histogramName, metric.WithDescription("Latency of the providers requests"))

	if err != nil {
		panic(err)
	}

	attr := kv{
		attribute.String("url", url),
	}

	histo.Record(context.Background(), latency, metric.WithAttributes(attr...))
}
