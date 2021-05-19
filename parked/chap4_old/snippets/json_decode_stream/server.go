package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type logLine struct {
	UserIP string `json:"user_ip"`
	Event  string `json:"event"`
}

func unmarshalHandler(w http.ResponseWriter, r *http.Request) {

	var l logLine
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &l)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v: %v\n", l.UserIP, l.Event)

}

func decodeHandler(w http.ResponseWriter, r *http.Request) {

	dec := json.NewDecoder(r.Body)
	var l logLine

	for {
		err := dec.Decode(&l)
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Println(l.UserIP, l.Event)
	}
	fmt.Fprintf(w, "OK")
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/unmarshal", unmarshalHandler)
	mux.HandleFunc("/decode", decodeHandler)

	http.ListenAndServe(":8080", mux)
}
