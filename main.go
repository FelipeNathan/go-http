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
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

var (
	insecure bool
	url      string
	method   string
	certPath string
	serve    bool
	client   *httpclient.HttpClient
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

	var err error
	client, err = httpclient.NewHttpClient(insecure, certPath)
	if err != nil {
		panic(err)
	}

	metric.Config()
	defer metric.Shutdown()

	if serve {
		httpServer()
	} else {
		res := makeRequest()
		fmt.Println(res)
	}
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

func httpServer() {
	r := chi.NewRouter()
	// useLogger(r)
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

func useLogger(r *chi.Mux) {
	logger := httplog.NewLogger("go-http", httplog.Options{
		JSON: true,
	})
	r.Use(httplog.RequestLogger(logger))
}
