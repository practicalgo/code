package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, req *http.Request, logger *log.Logger) (int, error) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "ok")
	return 200, nil
}

func deepCheckHandler(w http.ResponseWriter, req *http.Request, logger *log.Logger) (int, error) {
	logger.Print("Handling deepcheck")
	_, ok := req.URL.Query()["error"]
	if ok {
		return 500, errors.New("Error while running deepcheck")
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "deepcheck_ok")
	return 200, nil
}
