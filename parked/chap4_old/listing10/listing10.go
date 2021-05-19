package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	code int
}

func (mw *CustomResponseWriter) WriteHeader(code int) {
	mw.code = code
	mw.ResponseWriter.WriteHeader(code)
}

type app struct {
	logger  *log.Logger
	handler func(rw http.ResponseWriter, req *http.Request, logger *log.Logger) (int, error)
}

func (a *app) ServeHTTP(r http.ResponseWriter, req *http.Request) {
	httpStatus, err := a.handler(r, req, a.logger)
	if err != nil {
		http.Error(r, err.Error(), httpStatus)
	}
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := chi.NewRouter()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	wrappedMux := setupHandlers(mux, logger)

	log.Printf("Server attempting to listen on: %s\n", listenAddr)
	err := http.ListenAndServe(listenAddr, wrappedMux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}

func setupHandlers(mux *chi.Mux, logger *log.Logger) http.Handler {
	mux.Method("GET", "/healthcheck", &app{logger: logger, handler: healthCheckHandler})
	mux.Method("GET", "/deepcheck", &app{logger: logger, handler: deepCheckHandler})
	mux.Method("GET", "/", &app{logger: logger, handler: indexGetHandler})
	mux.Method("POST", "/", &app{logger: logger, handler: indexPostHandler})

	adminRouter := chi.NewRouter()
	adminRouter.Method("GET", "/login", &app{logger: logger, handler: adminLoginHandler})
	adminRouter.Method("GET", "/{adminId:[0-9]+}", &app{logger: logger, handler: getAdminHandler})
	adminRouter.Method("GET", "/{adminName:[a-z]+}", &app{logger: logger, handler: getAdminHandler})

	mux.Mount("/admin", adminRouter)

	return logRequestsMiddleware(
		logger,
		panicRecoveryMiddleware(
			logger,
			stripTrailingSlashMiddleware(logger, mux)))
}

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
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hello, world! POST handler")
	return 200, nil
}

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

func logRequestsMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myrw := &CustomResponseWriter{ResponseWriter: w, code: -1}
		next.ServeHTTP(myrw, r)
		if myrw.code == -1 {
			panic(fmt.Sprintf("HTTP response status not set in handler: %#v", next))
		}
		logger.Printf("method=%s path=%s status=%d\n", r.Method, r.RequestURI, myrw.code)
	})
}

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

func panicRecoveryMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				logger.Printf("stacktrace=%s\n", string(debug.Stack()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func stripTrailingSlashMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}
