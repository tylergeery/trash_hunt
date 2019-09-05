package connection

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/model"
	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

// Move is a move that a client wants to make
type Move struct {
	pos      game.Pos
	matchID  int64
	playerID int64
}

func NewMove(pos game.Pos, matchID, playerID int64) Move {
	return Move{
		pos:      pos,
		matchID:  matchID,
		playerID: playerID,
	}
}

// Client holds the client player information
type Client struct {
	conn          *net.TCPConn
	matchID       int64
	moveChan      chan Move
	notifications chan int
	player        *game.Player
}

// NewClient returns a new active client
func NewClient(conn *net.TCPConn) *Client {
	return &Client{
		conn:          conn,
		notifications: make(chan int, 5),
	}
}

// ValidateUser ensures that we have a valid user player
func (c *Client) ValidateUser() error {
	msg := make([]byte, 250)
	userToken, err := c.gatherInput(msg)
	if err != nil {
		fmt.Printf("Client: error getting user token: %s\n", err)
		return err
	}

	playerID, err := auth.GetPlayerIDFromAccessToken(strings.TrimRight(userToken, " \n\t"))
	fmt.Printf("Client: token %s\n", userToken)
	if err != nil {
		fmt.Printf("Client: could not validate token, received err: %s\n", err)
		return err
	}

	player := model.PlayerGetByID(playerID)
	if player.ID == 0 {
		fmt.Println("Client: could not find player in token")
		c.respond(err.Error())
		return fmt.Errorf("User with id (%d) could not be found", playerID)
	}

	fmt.Println("Client: Got player %s", player)
	c.player = game.NewPlayer(playerID)

	return nil
}

// GetMove collects a move from the client
func (c *Client) GetMove() (string, error) {
	msg := make([]byte, 2)
	move, err := c.gatherInput(msg)

	if err != nil {
		return "", err
	}

	return string(move), nil
}

func (c *Client) gatherInput(input []byte) (s string, err error) {
	_, err = c.conn.Read(input)
	if err != nil {
		return
	}

	s = string(bytes.TrimRight(input, "\x00"))

	return
}

// Respond sends output to the client
func (c *Client) processGame() {
	for {
		fmt.Println("Client: processing game...")
		move, err := c.gatherInput(make([]byte, 5))
		if err != nil {
			fmt.Printf("Client: ending game, %s\n", err)
			return
		}

		switch move {
		case "l":
			pos := game.Pos{X: c.player.Pos.X - 1, Y: c.player.Pos.Y}
			move := NewMove(pos, c.matchID, c.player.ID)
			c.moveChan <- move
		}
	}
}

// Respond sends output to the client
func (c *Client) respond(message string) error {
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Client: error sending response: %s\n", message)
	}

	return err
}

// WaitForStart holds until the game begins for clients
func (c *Client) WaitForStart() {
	fmt.Println("Client: waiting for start")

	err := c.respond("Status: Pending")
	if err != nil {
		return
	}

	select {
	case move := <-c.notifications:
		if move == 1 {
			// game start move
			c.processGame()

			return
		}
		fmt.Printf("Client: Received move before game started (%s)", move)
	}
}
