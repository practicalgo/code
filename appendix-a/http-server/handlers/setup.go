package handlers

import (
	"net/http"

	"github.com/practicalgo/code/appendix-a/http-server/config"
	"github.com/practicalgo/code/appendix-a/http-server/types"
)

func SetupHandlers(mux *http.ServeMux, config *config.AppConfig) {
	mux.Handle(
		"/api/packages",
		&types.App{Config: config, Handler: packageHandler},
	)
}
