package controllers

import (
	"net/http"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/tylergeery/trash_hunt/src/api_server/requests"
	"github.com/tylergeery/trash_hunt/src/api_server/responses"
	"github.com/tylergeery/trash_hunt/src/auth"
	"github.com/tylergeery/trash_hunt/src/game"
)

// PlayerCreate - Register/Signup a player
func PlayerCreate(c *routing.Context) error {
	// gather player info
	var req requests.PlayerCreateRequest
	err := c.Read(&req)
	if err != nil {
		return routing.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// insert new player into DB
	player, err := game.PlayerRegister(req.Email, req.Pw, req.Name, req.FacebookID)
	if err != nil {
		return routing.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// create a (permanent) jwt token for player
	token, err := auth.CreateToken(player)
	if err != nil {
		return routing.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// return player to the client
	resp := responses.PlayerCreateResponse{
		Token: token,
	}

	return c.Write(resp)
}

// PlayerLogin - Login the current player (get their auth token)
func PlayerLogin(c *routing.Context) error {
	var req requests.PlayerLoginRequest
	err := c.Read(&req)
	if err != nil {
		return err
	}

	// TODO: check if user has been rate limited
	p, err := game.PlayerLogin(req.Email, req.Pw)
	if err != nil {
		// TODO: add rate limiting for failures on email/ip
		return routing.NewHTTPError(http.StatusBadRequest, "Invalid credentials")
	}

	token, err := auth.CreateToken(p)
	if err != nil {
		return routing.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Return User w/ token
	resp := responses.PlayerLoginResponse{
		Player: p,
		Token:  token,
	}

	return c.Write(resp)
}

// PlayerUpdate - Edit a player
func PlayerUpdate(c *routing.Context) error {
	return nil
}

// PlayerDelete - Delete a player
func PlayerDelete(c *routing.Context) error {
	return nil
}

// PlayerQuery - Get information for a given player
func PlayerQuery(c *routing.Context) error {
	return nil
}

// PlayerResetPassword performs the password reset operation for a user
func PlayerResetPassword(c *routing.Context) error {
	return nil
}
