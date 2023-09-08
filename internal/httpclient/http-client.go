package httpclient

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/FelipeNathan/go-http/internal/apm/jaeger"
	"github.com/FelipeNathan/go-http/internal/httpclient/config"
	"github.com/FelipeNathan/go-http/internal/metric"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type HttpClient struct {
	*http.Client
}

func NewHttpClient(insecure bool, certPath string) (client *HttpClient, err error) {

	tc, err := config.TransportConfig(insecure, certPath)
	httpClient := &http.Client{Transport: tc}

	if err != nil {
		return nil, err
	}

	client = &HttpClient{
		Client: httpClient,
	}

	return
}

// Em construção
func (c *HttpClient) Post(ctx context.Context, url string) string {
	return c.doRequest(ctx, http.MethodPost, url)
}

func (c *HttpClient) Get(ctx context.Context, url string) string {
	return c.doRequest(ctx, http.MethodGet, url)
}

func (c *HttpClient) doRequest(ctx context.Context, method string, url string) string {

	_, span := jaeger.Trace(ctx, "makeRequest", attribute.String("url", url))
	defer span.End()

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	res, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	saveMetrics(now, url, res, span)

	body, _ := io.ReadAll(res.Body)
	return string(body)
}

func saveMetrics(start time.Time, url string, res *http.Response, span trace.Span) {
	latency := time.Since(start)
	metric.Count(url, res.StatusCode)
	metric.Gauge(url, latency.Milliseconds())
	metric.Histogram(url, latency.Milliseconds())
	span.SetAttributes(attribute.Int64("latency", latency.Milliseconds()))
}
