package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/practicalgolang/code/chap2/listing18/middleware"
	"github.com/practicalgolang/code/chap2/listing18/types"
)

func SetupHandlers(mux *chi.Mux, logger *log.Logger) http.Handler {
	mux.Method("GET", "/healthcheck", &types.App{Logger: logger, Handler: healthCheckHandler})
	mux.Method("GET", "/deepcheck", &types.App{Logger: logger, Handler: deepCheckHandler})
	mux.Method("GET", "/", &types.App{Logger: logger, Handler: indexGetHandler})
	mux.With(middleware.AuthMiddleware).Method("POST", "/", &types.App{Logger: logger, Handler: indexPostHandler})

	adminRouter := chi.NewRouter()
	adminRouter.Method("GET", "/login", &types.App{Logger: logger, Handler: adminLoginHandler})
	adminRouter.Method("GET", "/{adminId:[0-9]+}", &types.App{Logger: logger, Handler: getAdminHandler})
	adminRouter.Method("GET", "/{adminName:[a-z]+}", &types.App{Logger: logger, Handler: getAdminHandler})

	mux.Mount("/admin", adminRouter)

	return middleware.LogRequestsMiddleware(logger,
		http.TimeoutHandler(middleware.PanicRecoveryMiddleware(logger,
			middleware.StripTrailingSlashMiddleware(logger, mux)),
			30*time.Second,
			"Your request couldn't be processed by us",
		))
}
