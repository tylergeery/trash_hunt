package connection

import (
	"encoding/json"
	"fmt"

	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

var (
	messageStatusPending   = 0
	messageError           = 1
	messageInitGame        = 2
	messageUpdateGameState = 3
	messageEndGame         = 4

	eventStartGame = 1
	eventEndGame   = 2
)

// GameMessage is event sent to client through connection
type GameMessage struct {
	Event int    `json:"event"`
	Data  string `json:"data"`
}

// NewGameMessage creates a new event to send to client
func NewGameMessage(event int, data string) GameMessage {
	return GameMessage{event, data}
}

// ToBytes gets a game message as byte slice to send to client
func (gm GameMessage) ToBytes() []byte {
	bytes, err := json.Marshal(gm)
	if err != nil {
		fmt.Println("GameMessage Marshal error:", err.Error())
	}

	return bytes
}

// GameSetUp is the obect for different game settings
type GameSetUp struct {
	UserToken  string `json:"user_token"`
	Opponent   int64  `json:"opponent"`
	Difficulty string `json:"difficulty"`
}

// Move is a move that a client wants to make
type Move struct {
	pos      game.Pos
	matchID  int64
	playerID int64
}

// NewMove returns a new game move object
func NewMove(pos game.Pos, matchID, playerID int64) Move {
	return Move{
		pos:      pos,
		matchID:  matchID,
		playerID: playerID,
	}
}
