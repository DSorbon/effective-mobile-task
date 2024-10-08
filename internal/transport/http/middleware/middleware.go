package middleware

import (
	"net/http"
	"time"

	"github.com/DSorbon/effective-mobile-task/pkg/logger"
)

func LogMiddleware(next http.Handler) http.Handler {

	res := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()

		next.ServeHTTP(w, r)

		logger.Infof("request: path-%v, method-%v, duration-%v", r.URL.Path, r.Method, time.Since(now))
	})

	return res
}
