package middleware

import (
	"net/http"
	"strings"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/tylergeery/trash_hunt/auth"
)

type key string

const (
	// PlayerIDKey for storing player_id in request context
	PlayerIDKey key = "player_id"
)

func validateToken(c *routing.Context) error {
	var err error

	bearer := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(bearer, "Bearer ")
	playerID, err := auth.GetPlayerIDFromAccessToken(token)

	if err != nil {
		return routing.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	c.Set("PlayerID", playerID)

	return nil
}
