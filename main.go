package main

import (
	"flag"
	"fmt"
	"strings"

	httpclient "github.com/FelipeNathan/go-http/http-client"
)

var (
	insecure bool
	url      string
	method   string
	certPath string
)

func main() {

	flag.BoolVar(&insecure, "insecure", false, "Ignore TLS validation")
	flag.StringVar(&url, "url", "", "Url to make the request")
	flag.StringVar(&method, "method", "GET", "Http method")
	flag.StringVar(&certPath, "certPath", "./certs/", "Path of pem certificates to trust")
	flag.Parse()

	if !strings.HasSuffix(certPath, "/") {
		certPath += "/"
	}

	makeRequest()
}

func makeRequest() {
	client, err := httpclient.NewHttpClient(insecure, certPath)
	if err != nil {
		panic(err)
	}

	var res string
	switch method {
	case "POST":
		res = client.Post(url)
	default:
		res = client.Get(url)
	}

	fmt.Println(res)
}
