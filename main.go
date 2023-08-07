package main

import (
	"flag"
	"fmt"

	httpclient "github.com/FelipeNathan/go-http/http-client"
)

var (
	insecure bool
	url      string
	method   string
	client   *httpclient.HttpClient
)

func main() {

	flag.BoolVar(&insecure, "insecure", false, "Ignore TLS validation")
	flag.StringVar(&url, "url", "", "Url to make the request")
	flag.StringVar(&method, "method", "GET", "Http method")
	flag.Parse()

	makeRequest()
}

func makeRequest() {
	var err error
	client, err = httpclient.NewHttpClient(insecure)
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
