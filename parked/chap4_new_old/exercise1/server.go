package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

func apiHandler(w http.ResponseWriter, req *http.Request) {
	type Inventory struct {
		Material string
		Count    uint
	}
	sweaters := Inventory{"wool", 17}
	tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, sweaters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")
}

func setupHandlers() {
	http.HandleFunc("/healthcheck", healthCheckHandler)
	http.HandleFunc("/api", apiHandler)

}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	setupHandlers()

	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}
