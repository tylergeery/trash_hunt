package connection

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/model"
	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

// Client holds the client player information
type Client struct {
	conn          Connection
	matchID       int64
	moveChan      chan Move // For sending moves to manager
	notifications chan int  // For receiving events from manager/arena
	player        *game.Player
	preferences   GameSetUp
}

// NewClient returns a new active client
func NewClient(conn Connection) *Client {
	return &Client{
		conn:          conn,
		notifications: make(chan int, 5),
	}
}

// SetUpUser ensures that we have a valid user player
func (c *Client) SetUpUser() error {
	var gameSetUp GameSetUp
	msg := make([]byte, 500)
	settings, err := c.conn.gatherInput(msg)
	if err != nil {
		fmt.Printf("Client: error getting user game options: %s\n", err)
		return err
	}

	err = json.Unmarshal([]byte(settings), &gameSetUp)
	if err != nil {
		fmt.Printf("Client: error unmarshaling user game set up from (%s): %s\n", settings, err)
		return err
	}

	c.preferences = gameSetUp
	playerID, err := auth.GetPlayerIDFromAccessToken(strings.TrimRight(gameSetUp.UserToken, " \n\t"))
	fmt.Printf("Client: token %s\n", gameSetUp.UserToken)
	if err != nil {
		fmt.Printf("Client: could not validate token, received err: %s\n", err)
		return err
	}

	player := model.PlayerGetByID(playerID)
	if player.ID == 0 {
		fmt.Println("Client: could not find player in token")
		c.conn.respond(NewGameMessage(messageError, err.Error()))
		return fmt.Errorf("User with id (%d) could not be found", playerID)
	}

	fmt.Printf("Client: Got player (%d)\n", player.ID)
	c.player = game.NewPlayer(playerID)

	return nil
}

// GetMove collects a move from the client
func (c *Client) GetMove() (string, error) {
	msg := make([]byte, 2)
	move, err := c.conn.gatherInput(msg)

	if err != nil {
		return "", err
	}

	return string(move), nil
}

func (c *Client) processGame() {
	for {
		var pos game.Pos
		fmt.Println("Client: processing game...")
		clientMove, err := c.conn.gatherInput(make([]byte, 5))
		if err != nil {
			fmt.Printf("Client: ending game, %s\n", err)
			return
		}

		switch clientMove {
		case "l":
			pos = game.Pos{X: c.player.Pos.X - 1, Y: c.player.Pos.Y}
		case "r":
			pos = game.Pos{X: c.player.Pos.X + 1, Y: c.player.Pos.Y}
		case "u":
			pos = game.Pos{X: c.player.Pos.X, Y: c.player.Pos.Y - 1}
		case "d":
			pos = game.Pos{X: c.player.Pos.X, Y: c.player.Pos.Y + 1}
		default:
			fmt.Printf("Client: unknown move: %s", clientMove)
			continue
		}

		move := NewMove(pos, c.matchID, c.player.ID)
		c.moveChan <- move
	}
}

// WaitForStart holds until the game begins for clients
func (c *Client) WaitForStart() {
	fmt.Println("Client: waiting for start")

	msg := NewGameMessage(messageStatusPending, "Status: Pending")
	err := c.conn.respond(msg)
	if err != nil {
		return
	}

	defer func() {
		fmt.Println("Client: game finished")
	}()

	select {
	case event := <-c.notifications:
		if event == eventStartGame {
			c.processGame()
		}

		fmt.Printf("Client: Received move before game started (%s)\n", string(event))
	}
}
