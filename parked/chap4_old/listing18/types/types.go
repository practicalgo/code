package types

import (
	"log"
	"net/http"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	Code         int
	BytesWritten int
	wroteHeader  bool
}

func (mw *CustomResponseWriter) WriteHeader(code int) {
	mw.Code = code
	mw.ResponseWriter.WriteHeader(code)
	mw.wroteHeader = true
}

func (mw *CustomResponseWriter) Write(p []byte) (int, error) {
	if !mw.wroteHeader {
		mw.WriteHeader(http.StatusOK)
	}
	mw.BytesWritten += len(p)
	return mw.ResponseWriter.Write(p)
}

type App struct {
	Logger  *log.Logger
	Handler func(rw http.ResponseWriter, req *http.Request, logger *log.Logger) (int, error)
}

func (a *App) ServeHTTP(r http.ResponseWriter, req *http.Request) {
	httpStatus, err := a.Handler(r, req, a.Logger)
	if err != nil {
		http.Error(r, err.Error(), httpStatus)
	}
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"alpha,min=3,max=15"`
	Password string `json:"password" validate:"alphanum,min=6,max=15"`
}

type UserLoginResponse struct {
	UserId   uint8  `json:"id"`
	Username string `json:"username"`
}
