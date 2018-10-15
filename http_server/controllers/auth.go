package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tylergeery/trash_hunt/auth"
)

// Auth - Create a new auth token from user key
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		defer func(begin time.Time) {
			log.Printf("Authorized %t\n", err == nil)
		}(time.Now())

		bearer := r.Header.Get("Authorization")
		token := strings.TrimPrefix(bearer, "Bearer ")
		playerID, err := auth.GetPlayerIDFromAccessToken(token)

		if err == nil {
			ctx := context.WithValue(r.Context(), "player_id", playerID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// TODO: reject request?
			next.ServeHTTP(w, r)
		}
	})
}
