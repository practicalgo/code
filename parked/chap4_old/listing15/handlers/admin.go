package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func adminLoginHandler(w http.ResponseWriter, req *http.Request, logger *log.Logger) (int, error) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "admin_login")
	return 200, nil
}

func getAdminHandler(w http.ResponseWriter, req *http.Request, logger *log.Logger) (int, error) {
	var msg string
	if adminId := chi.URLParam(req, "adminId"); adminId != "" {
		msg = fmt.Sprintf("admin_%v", adminId)
	} else if adminName := chi.URLParam(req, "adminName"); adminName != "" {
		msg = fmt.Sprintf("admin_%v", adminName)
	} else {
		return 400, nil
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, msg)
	return 200, nil
}
