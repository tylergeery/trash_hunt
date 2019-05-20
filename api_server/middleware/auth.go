package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/tylergeery/trash_hunt/auth"
)

type key string

const (
	// PlayerIDKey for storing player_id in request context
	PlayerIDKey key = "player_id"
)

func validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		bearer := r.Header.Get("Authorization")
		token := strings.TrimPrefix(bearer, "Bearer ")
		playerID, err := auth.GetPlayerIDFromAccessToken(token)

		if err == nil {
			ctx := context.WithValue(r.Context(), PlayerIDKey, playerID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// TODO: reject request?
			next.ServeHTTP(w, r)
		}
	})
}
