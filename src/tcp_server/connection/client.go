package connection

import (
	"fmt"
	"net"
	"strings"
	"time"

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
	msg := make([]byte, 250)
	err := c.gatherInput(msg)
	if err != nil {
		fmt.Printf("Client: error getting user token: %s\n", err)
		return err
	}

	userToken := strings.TrimRight(string(msg), " \n\t")
	playerID, err := auth.GetPlayerIDFromAccessToken(userToken)
	fmt.Printf("Client: token %s after\n", userToken)
	fmt.Printf("Client: token %s\n", msg[:42])
	fmt.Printf("Client: token %s\n", msg[42:])
	if err != nil {
		fmt.Printf("Client: could not validatse token, received err: %s\n", err)
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
	for {
		fmt.Println("Client: processing game...")
		time.Sleep(200 * time.Millisecond)
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

	// TODO: Listen for Manager to say game begins
	err := c.respond("starting")
	if err != nil {
		return
	}

	c.processGame()
}
