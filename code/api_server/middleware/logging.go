package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-ozzo/ozzo-routing"
)

func logRequest() func(c *routing.Context) error {
	return func(begin time.Time) func(c *routing.Context) error {
		return func(c *routing.Context) error {
			r := c.Request
			log.Printf("%s %s %s - %dμs\n", r.Method, r.URL.Path, getToken(r), time.Since(begin).Nanoseconds()/1000)

			return nil
		}
	}(time.Now())
}

func getToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}
