package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/game"
)

type key string

const (
	// PlayerIDKey for storing player_id in request context
	PlayerIDKey key = "player_id"
)

// CreateAuthToken generates a new auth token for use
func CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Get player key
	key := r.Form.Get("key")

	// Look up player
	player := game.PlayerGetByToken(key[0])

	// Create temp auth token
	token, err := auth.CreateToken(player)

	// Send back to client
	if err == nil {
		response, err := json.Marshal()
		if err == nil {
			w.Write(response)
		}
	}

	// Set Error Response
}
