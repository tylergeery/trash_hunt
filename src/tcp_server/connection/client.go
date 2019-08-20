package connection

import (
	"net"

	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

// Client holds the client player information
type Client struct {
	conn   net.Conn
	player *game.Player
}

// NewClient returns a new active client
func NewClient(conn net.Conn, player *game.Player) *Client {
	return &Client{
		conn:   conn,
		player: player,
	}
}

// GatherMove collects a move from the client
func (c *Client) GatherMove() {

}

// Respond sends output to the client
func (c *Client) Respond() {

}
