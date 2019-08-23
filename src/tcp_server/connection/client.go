package connection

import (
	"fmt"
	"net"

	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/model"
	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

// Client holds the client player information
type Client struct {
	conn    net.Conn
	matchID int64
	player  *game.Player
}

// NewClient returns a new active client
func NewClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

// ValidateUser ensures that we have a valid user player
func (c *Client) ValidateUser() error {
	userToken := make([]byte, 180)
	err := c.gatherInput(userToken)
	if err != nil {
		return err
	}

	playerID, err := auth.GetPlayerIDFromAccessToken(string(userToken))
	if err != nil {
		c.respond(err.Error())
		return err
	}

	player := model.PlayerGetByID(playerID)
	if player.ID == 0 {
		c.respond(err.Error())
		return fmt.Errorf("User with id (%d) could not be found", playerID)
	}

	c.player = game.NewPlayer(playerID)

	return nil
}

// GetMove collects a move from the client
func (c *Client) GetMove() (string, error) {
	move := make([]byte, 2)
	err := c.gatherInput(move)

	if err != nil {
		return "", err
	}

	return string(move), nil
}

func (c *Client) gatherInput(input []byte) error {
	_, err := c.conn.Read(input)

	if err != nil {
		return err
	}

	return nil
}

// Respond sends output to the client
func (c *Client) processGame() {

}

// Respond sends output to the client
func (c *Client) respond(message string) error {
	_, err := c.conn.Write([]byte(message))

	return err
}

// WaitForStart holds until the game begins for clients
func (c *Client) WaitForStart() {
	// TODO: Listen for Manager to say game begins
	c.respond("starting")
	c.processGame()
}
