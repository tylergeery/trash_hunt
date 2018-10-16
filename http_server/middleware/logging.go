package middleware

import (
	"log"
	"net/http"
	"time"
)

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(begin time.Time) {
			log.Printf("%s %s - %dms\n", r.Method(), r.URL.Path, time.Since(begin).Milliseconds())
		}(time.Now())

		next.ServeHTTP(w, r)
	})
}
