package httpclient

import (
	"io"
	"net/http"

	"github.com/FelipeNathan/go-http/httpclient/config"
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

	res, err := c.Do(req)

	if err != nil {
		panic(err)
	}

	body, _ := io.ReadAll(res.Body)
	return string(body)
}
