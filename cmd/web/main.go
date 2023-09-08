package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/FelipeNathan/go-http/internal/apm/jaeger"
	"github.com/FelipeNathan/go-http/internal/controller"
	"github.com/FelipeNathan/go-http/internal/metric"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	certPath string
)

func main() {

	flag.StringVar(&certPath, "certPath", "./assets/certs/", "Path of pem certificates to trust")
	flag.Parse()

	if !strings.HasSuffix(certPath, "/") {
		certPath += "/"
	}

	metric.Config()
	defer metric.Shutdown()

	jaeger.Config()
	defer jaeger.Shutdown()

	httpServer()
}

func httpServer() {
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(getLogger()))
	r.Use(middleware.Logger)

	controller.CertPath = certPath

	addRoute(r, http.MethodGet, "/", controller.Index)
	addRoute(r, http.MethodPost, "/", controller.Post)

	fmt.Print("Listening")
	http.ListenAndServe(":8080", r)
}

func getLogger() zerolog.Logger {
	return httplog.NewLogger("go-http", httplog.Options{
		JSON: false,
	})
}

func addRoute(r *chi.Mux, method, route string, handler func(w http.ResponseWriter, req *http.Request)) {
	r.Method(method, route, otelhttp.NewHandler(http.HandlerFunc(handler), route))
}
