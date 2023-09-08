package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/FelipeNathan/go-http/internal/httpclient"
	"github.com/FelipeNathan/go-http/internal/metric"
)

var (
	insecure bool
	url      string
	method   string
	certPath string
	client   *httpclient.HttpClient
)

func main() {

	flag.BoolVar(&insecure, "insecure", false, "Ignore TLS validation")
	flag.StringVar(&url, "url", "", "Url to make the request")
	flag.StringVar(&method, "method", "GET", "Http method")
	flag.StringVar(&certPath, "certPath", "./assets/certs/", "Path of pem certificates to trust")
	flag.Parse()

	if !strings.HasSuffix(certPath, "/") {
		certPath += "/"
	}

	var err error
	client, err = httpclient.NewHttpClient(insecure, certPath)
	if err != nil {
		panic(err)
	}

	metric.Config()
	defer metric.Shutdown()

	res := makeRequest()
	fmt.Println(res)
}

func makeRequest() string {
	var res string
	switch method {
	case "POST":
		res = client.Post(url)
	default:
		res = client.Get(url)
	}
	return res
}
