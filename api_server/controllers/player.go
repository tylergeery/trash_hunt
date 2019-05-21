package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/game"
)

// PlayerCreate - Register/Signup a player
func PlayerCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// gather player info
	email := r.Form.Get("email")
	pw := r.Form.Get("pw")
	name := r.Form.Get("name")
	facebookID := r.Form.Get("facebook_id")

	// insert new player into DB
	player, err := game.PlayerRegister(email, pw, name, facebookID)
	if err != nil {
		// TODO: generate error response
	}

	// create a (permanent) jwt token for player
	token, err := auth.CreateToken(player)
	if err != nil {
		// TODO: generate "logged out" response
	}

	// return player to the client
	json.Marshal(token)
}

// PlayerLogin - Login the current player (get their auth token)
func PlayerLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	_ = r.Form.Get("email")
	_ = r.Form.Get("pw")

	// Find user by email and password

	// TODO: add rate limiting for failures on email/ip

	// Return User w/ token
}

// PlayerUpdate - Edit a player
func PlayerUpdate(w http.ResponseWriter, r *http.Request) {

}

// PlayerDelete - Delete a player
func PlayerDelete(w http.ResponseWriter, r *http.Request) {

}

// PlayerQuery - Get information for a given player
func PlayerQuery(w http.ResponseWriter, r *http.Request) {

}
