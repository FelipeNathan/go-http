package controller

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/FelipeNathan/go-http/internal/httpclient"
)

var CertPath string

type payload struct {
	Url      string `json:"url"`
	Insecure bool   `json:"insecure"`
	Method   string `json:"method"`
}

func Index(w http.ResponseWriter, req *http.Request) {
	template, err := template.ParseFiles("./web/html/index.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	template.Execute(w, nil)
}

func Post(w http.ResponseWriter, req *http.Request) {

	p := &payload{}
	err := json.NewDecoder(req.Body).Decode(p)

	if err != nil {
		panic(err)
	}

	response := makeRequest(*p)

	w.Header().Add("Content-Type", "text/html")
	template.HTMLEscape(w, []byte(response))
}

func makeRequest(p payload) string {
	client, err := httpclient.NewHttpClient(p.Insecure, CertPath)
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
