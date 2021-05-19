package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/practicalgolang/code/chap2/listing16/types"
)

func LogRequestsMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myrw := &types.CustomResponseWriter{ResponseWriter: w, Code: -1}
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			logger.Printf("method=%s path=%s status=%d duration=%f bytes=%d\n", r.Method, r.RequestURI, myrw.Code, duration.Seconds(), myrw.BytesWritten)
		}()
		next.ServeHTTP(myrw, r)
		if myrw.Code == -1 {
			panic(fmt.Sprintf("HTTP response status not set in handler: %#v", next))
		}
	})
}

func PanicRecoveryMiddleware(logger *log.Logger, next http.Handler) http.Handler {
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

func StripTrailingSlashMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := types.UserLoginRequest{}
		resp := types.UserLoginResponse{}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err == nil {
			validate := validator.New()
			err = validate.Struct(u)
			if err == nil {
				resp.UserId = 1
				resp.Username = u.Username
			}
		}
		ctx := context.WithValue(r.Context(), "User", u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
