package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/tylergeery/trash_hunt/api_server/requests"
	"github.com/tylergeery/trash_hunt/api_server/responses"
	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/game"
)

type key string

// CreateAuthToken generates a new auth token for use
func CreateAuthToken(c *routing.Context) error {
	req := requests.NewCreateAuthTokenRequest(c)
	if err := req.Validate(); err != nil {
		return routing.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Look up player
	player := game.PlayerGetByToken(req.Key)

	// Create temp auth token
	token, err := auth.CreateToken(player)
	if err != nil {
		return routing.NewHTTPError(http.StatusBadRequest, "Invalid key supplied")
	}

	resp := responses.AuthTokenCreateResponse{Token: token}
	output, err := json.Marshal(resp)
	if err != nil {
		return routing.NewHTTPError(http.StatusBadRequest)
	}

	return c.Write(output)
}
