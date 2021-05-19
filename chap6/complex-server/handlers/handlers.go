package handlers

import (
	"fmt"
	"net/http"

	"github.com/practicalgo/code/chap6/complex-server/config"
)

type app struct {
	conf    config.AppConfig
	handler func(w http.ResponseWriter, r *http.Request, conf config.AppConfig)
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler(w, r, a.conf)
}

func apiHandler(w http.ResponseWriter, r *http.Request, conf config.AppConfig) {
	fmt.Fprintf(w, "Hello, world!")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, conf config.AppConfig) {
	if r.Method != "GET" {
		conf.Logger.Printf("error=\"Invalid request\" path=%s method=%s", r.URL.Path, r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "ok")
}

func panicHandler(w http.ResponseWriter, r *http.Request, conf config.AppConfig) {
	panic("I panicked")
}
