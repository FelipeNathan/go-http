package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/FelipeNathan/go-http/internal/httpclient"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

var (
	certPath string
)

type payload struct {
	Url      string `json:"url"`
	Insecure bool   `json:"insecure"`
	Method   string `json:"method"`
}

func main() {

	flag.StringVar(&certPath, "certPath", "./assets/certs/", "Path of pem certificates to trust")
	flag.Parse()

	if !strings.HasSuffix(certPath, "/") {
		certPath += "/"
	}

	// metric.Config()
	// defer metric.Shutdown()
	httpServer()
}

func makeRequest(p payload) string {
	client, err := httpclient.NewHttpClient(p.Insecure, certPath)
	if err != nil {
		panic(err)
	}

	var res string
	switch p.Method {
	case "POST":
		res = client.Post(p.Url)
	default:
		res = client.Get(p.Url)
	}
	return res
}

func httpServer() {
	r := chi.NewRouter()
	useLogger(r)
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		template, err := template.ParseFiles("./web/html/index.html")

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		template.Execute(w, nil)
	})

	r.Post("/", func(w http.ResponseWriter, req *http.Request) {

		p := &payload{}
		err := json.NewDecoder(req.Body).Decode(p)

		if err != nil {
			panic(err)
		}

		response := makeRequest(*p)

		w.Header().Add("Content-Type", "text/html")
		template.HTMLEscape(w, []byte(response))
	})

	fmt.Print("Listening")
	http.ListenAndServe(":8080", r)
}

func useLogger(r *chi.Mux) {
	logger := httplog.NewLogger("go-http", httplog.Options{
		JSON: false,
	})
	r.Use(httplog.RequestLogger(logger))
}
