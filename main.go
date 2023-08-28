package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/FelipeNathan/go-http/httpclient"
	"github.com/FelipeNathan/go-http/metric"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	insecure bool
	url      string
	method   string
	certPath string
	serve    bool
)

type payload struct {
	Url string `json:"url"`
}

func main() {

	flag.BoolVar(&insecure, "insecure", false, "Ignore TLS validation")
	flag.StringVar(&url, "url", "", "Url to make the request")
	flag.StringVar(&method, "method", "GET", "Http method")
	flag.StringVar(&certPath, "certPath", "./certs/", "Path of pem certificates to trust")
	flag.BoolVar(&serve, "serve", false, "Indicates if this should be served as a http server")
	flag.Parse()

	if !strings.HasSuffix(certPath, "/") {
		certPath += "/"
	}

	metric.Config()
	defer metric.Shutdown()

	if serve {
		httpServer()
	} else {
		makeRequest()
	}
}

func makeRequest() string {
	metric.Count()
	metric.Gauge()
	metric.Histogram()

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
	return res
}

func httpServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		template, _ := template.ParseFiles("./html/index.html")
		template.Execute(w, nil)
	})

	r.Post("/", func(w http.ResponseWriter, req *http.Request) {

		w.Header().Add("Content-Type", "text/html")
		p := &payload{}
		err := json.NewDecoder(req.Body).Decode(p)

		if err != nil {
			panic(err)
		}

		fmt.Print(p.Url)
		url = p.Url

		response := makeRequest()
		template.HTMLEscape(w, []byte(response))
	})

	http.ListenAndServe(":8080", r)
}
