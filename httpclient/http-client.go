package httpclient

import (
	"io"
	"net/http"
	"time"

	"github.com/FelipeNathan/go-http/httpclient/config"
	"github.com/FelipeNathan/go-http/metric"
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
func (c *HttpClient) Post(url string) string {
	return c.doRequest(http.MethodPost, url)
}

func (c *HttpClient) Get(url string) string {
	return c.doRequest(http.MethodGet, url)
}

func (c *HttpClient) doRequest(method string, url string) string {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	res, err := c.Do(req)
	latency := time.Since(now)

	if err != nil {
		panic(err)
	}

	metric.Count(url, res.StatusCode)
	metric.Gauge(url, latency.Milliseconds())
	metric.Histogram(url, latency.Milliseconds())

	body, _ := io.ReadAll(res.Body)
	return string(body)
}
