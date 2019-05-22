package controllers

import (
	"github.com/go-ozzo/ozzo-routing"
	"github.com/tylergeery/trash_hunt/api_server/requests"
	"github.com/tylergeery/trash_hunt/api_server/responses"
	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/game"
)

// PlayerCreate - Register/Signup a player
func PlayerCreate(c *routing.Context) error {
	// gather player info
	var req requests.PlayerCreateRequest
	err := c.Read(&req)
	if err != nil {
		return err
	}

	// insert new player into DB
	player, err := game.PlayerRegister(req.Email, req.Pw, req.Name, req.FacebookID)
	if err != nil {
		return err
	}

	// create a (permanent) jwt token for player
	token, err := auth.CreateToken(player)
	if err != nil {
		return err
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

	// Find user by email and password

	// TODO: add rate limiting for failures on email/ip

	// Return User w/ token
	return nil
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
