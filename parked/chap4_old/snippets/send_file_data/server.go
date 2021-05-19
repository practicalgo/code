package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func downloadFileHandler(w http.ResponseWriter, r *http.Request) {

	v := r.URL.Query()
	fileName := v.Get("fileName")
	if len(fileName) == 0 {
		http.Error(w, "fileName query parameter not specified", http.StatusBadRequest)
		return
	}

	f, err := os.Open(fileName)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	f.Seek(0, 0)
	contentType := http.DetectContentType(buffer)

	log.Println(contentType)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))

	io.Copy(w, f)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/download", downloadFileHandler)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}

}
