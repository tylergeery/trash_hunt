package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/tylergeery/pick-up-api/utils/response"
	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/game"
)

type key string

// CreateAuthToken generates a new auth token for use
func CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Get player key
	key := r.Form.Get("key")

	if key == "" {
		response.Fail(w, http.StatusBadRequest, "Invalid key supplied")
		return
	}

	// Look up player
	player := game.PlayerGetByToken(string(key[0]))

	// Create temp auth token
	token, err := auth.CreateToken(player)

	// Send back to client
	if err != nil {
		response.Fail(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := json.Marshal(map[string]string{"token": token})
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, resp)
}
