package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/practicalgo/code/chap6/log-response-status/config"
	"github.com/practicalgo/code/chap6/log-response-status/handlers"
	"github.com/practicalgo/code/chap6/log-response-status/middleware"
)

func setupServer(mux *http.ServeMux, w io.Writer) http.Handler {
	conf := config.InitConfig(w)

	handlers.Register(mux, conf)
	return middleware.RegisterMiddleware(mux, conf)
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	wrappedMux := setupServer(mux, os.Stdout)

	log.Fatal(http.ListenAndServe(listenAddr, wrappedMux))
}
