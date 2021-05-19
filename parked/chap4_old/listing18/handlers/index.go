package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/practicalgolang/code/chap2/listing18/types"
)

func indexGetHandler(w http.ResponseWriter, req *http.Request, logger *log.Logger) (int, error) {
	_, ok := req.URL.Query()["panic"]
	if ok {
		panic("Sorry, I couldn't process your request this time")
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hello, world!")
	return 200, nil
}

func indexPostHandler(w http.ResponseWriter, req *http.Request, logger *log.Logger) (int, error) {
	u, ok := req.Context().Value("User").(types.UserLoginResponse)
	if !ok || u.UserId == 0 {
		return 401, errors.New("Invalid user credentials")
	}

	w.WriteHeader(http.StatusOK)
	respJson, err := json.Marshal(u)
	if err != nil {
		return 500, err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
	return 200, nil
}
