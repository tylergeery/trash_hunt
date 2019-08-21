package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/tylergeery/trash_hunt/api_server/requests"
	"github.com/tylergeery/trash_hunt/api_server/responses"
	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/model"
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
	player, err := model.PlayerRegister(req.Email, req.Pw, req.Name, req.FacebookID)
	if err != nil {
		return routing.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// create a (permanent) jwt token for player
	token, err := auth.CreateToken(player)
	if err != nil {
		return routing.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// return player to the client
	resp := responses.PlayerLoginResponse{
		Player: player,
		Token:  token,
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
	p, err := model.PlayerLogin(req.Email, req.Pw)
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
	var req requests.PlayerUpdateRequest
	err := c.Read(&req)
	if err != nil {
		return err
	}

	authID := c.Get("PlayerID").(int64)
	player := model.PlayerGetByID(authID)

	player.Email = req.Email
	player.Name = req.Name
	player.FacebookID = req.FacebookID
	err = player.Update()

	if err != nil {
		return err
	}

	return c.Write(model.PlayerGetByID(authID))
}

// PlayerDelete - Delete a player
func PlayerDelete(c *routing.Context) error {
	authID := c.Get("PlayerID").(int64)
	player := model.PlayerGetByID(authID)
	player.Status = model.PlayerStatusRemoved
	c.Response.WriteHeader(http.StatusNoContent)

	return player.Update()
}

// PlayerQuery - Get information for a given player
func PlayerQuery(c *routing.Context) error {
	id := c.Param("id")
	playerID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	player := model.PlayerGetByID(int64(playerID))
	authID := c.Get("PlayerID").(int64)

	if authID == int64(playerID) {
		return c.Write(player)
	}

	return c.Write(player.ToPublicProfile())
}

// PlayerResetPassword performs the password reset operation for a user
func PlayerResetPassword(c *routing.Context) error {
	// TODO: need to set up email
	return nil
}
