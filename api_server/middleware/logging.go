package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(begin time.Time) {
			log.Printf("%s %s %s - %dμs\n", r.Method, r.URL.Path, getToken(r), time.Since(begin).Nanoseconds()/1000)
		}(time.Now())

		next.ServeHTTP(w, r)
	})
}

func getToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}
