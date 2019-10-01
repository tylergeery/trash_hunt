package connection

import (
	"encoding/json"
	"fmt"

	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

// Arena holds all the information about clients and the active game
type Arena struct {
	state   *game.State
	clients []*Client
}

// NewArena sets up all the infrastructure for gameplay
func NewArena(p1, p2 *game.Player, clients ...*Client) *Arena {
	return &Arena{
		state:   game.NewState(p1, p2),
		clients: clients,
	}
}

func (a *Arena) start(matchID int64, moveChan chan Move) {
	// TODO: DO better than this, this will race
	for i := range a.clients {
		if a.clients[i] == nil {
			continue
		}

		a.clients[i].matchID = matchID
		a.clients[i].moveChan = moveChan
	}

	a.notifyClients(1)
	a.sendInitialState()
}

func (a *Arena) notifyClients(move int) {
	for i := range a.clients {
		if a.clients[i] == nil {
			continue
		}

		a.clients[i].notifications <- move
	}
}

func (a *Arena) sendInitialState() {
	message, _ := json.Marshal(a.state)
	gameState := string(message)

	for i := range a.clients {
		a.clients[i].respond(gameState)
	}
}

func (a *Arena) sendPositions() {
	message, _ := json.Marshal(a.state.Players)
	fmt.Printf("Arena: sending positions: %s", message)
	positions := string(message)

	for i := range a.clients {
		if a.clients[i] == nil {
			continue
		}

		a.clients[i].respond(positions)
	}
}
