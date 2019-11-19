package connection

import (
	"encoding/json"

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
		a.clients[i].matchID = matchID
		a.clients[i].moveChan = moveChan
	}

	a.notifyClients(moveStartGame)
	a.sendInitialState()
}

func (a *Arena) notifyClients(move int) {
	for i := range a.clients {
		a.clients[i].notifications <- move
	}
}

func (a *Arena) sendInitialState() {
	message, _ := json.Marshal(a.state)
	gameState := string(message)

	for i := range a.clients {
		msg := NewGameMessage(messageInitGame, gameState)
		a.clients[i].conn.respond(msg)
	}
}

func (a *Arena) moveUser(playerID int64, nextPos game.Pos) {
	moved := a.state.MoveUser(playerID, nextPos)
	if !moved {
		return
	}

	for id := range a.state.Players {
		if id == playerID {
			continue
		}
		if id < 0 {
			// move computer player (as response)
			nextMove := a.state.Players[id].Solver.GetMove(id, a.state)
			_ = a.state.MoveUser(id, nextMove)
		}
	}
}

func (a *Arena) sendPositions() {
	message, _ := json.Marshal(a.state.Players)
	positions := string(message)

	for i := range a.clients {
		msg := NewGameMessage(messageUpdateGameState, positions)
		a.clients[i].conn.respond(msg)
	}
}
