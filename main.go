package main

import (
	"fmt"
	"net/http"

	httpclient "github.com/FelipeNathan/go-http/http-client"
	"github.com/FelipeNathan/go-http/metric"
)

var (
	insecure bool
	url      string
	method   string
	certPath string
)

func main() {

	metric.Config()
	defer metric.Shutdown()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		metric.Count()
		metric.Gauge()
		metric.Histogram()
		w.Write([]byte("Working"))
	})

	http.ListenAndServe(":8080", nil)

	// flag.BoolVar(&insecure, "insecure", false, "Ignore TLS validation")
	// flag.StringVar(&url, "url", "", "Url to make the request")
	// flag.StringVar(&method, "method", "GET", "Http method")
	// flag.StringVar(&certPath, "certPath", "./certs/", "Path of pem certificates to trust")
	// flag.Parse()

	// if !strings.HasSuffix(certPath, "/") {
	// 	certPath += "/"
	// }

	// makeRequest()
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
