package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type logLine struct {
	UserIP string `json:"user_ip"`
	Event  string `json:"event"`
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)

	var e *json.UnmarshalTypeError

	for {
		var l logLine
		err := dec.Decode(&l)
		if err == io.EOF {
			break
		}
		if errors.As(err, &e) {
			log.Println(err)
			continue
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(l.UserIP, l.Event)
	}
	fmt.Fprintf(w, "OK")
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/decode", decodeHandler)

	http.ListenAndServe(":8080", mux)
}
