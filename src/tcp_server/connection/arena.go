package connection

import (
	"encoding/json"

	"github.com/tylergeery/trash_hunt/model"
	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

// Arena holds all the information about clients and the active game
type Arena struct {
	state   *game.State
	match   *model.Game
	clients []*Client
}

// NewArena sets up all the infrastructure for gameplay
func NewArena(p1, p2 *game.Player, clients ...*Client) *Arena {
	return &Arena{
		state:   game.NewState(p1, p2),
		clients: clients,
	}
}

// HasWinner returns whether a winner has been acheived
func (a *Arena) HasWinner() bool {
	if a.state.GetWinner() != 0 {
		return true
	}

	return false
}

func (a *Arena) end() {
	if !a.HasWinner() {
		return
	}

	type Results struct {
		Winner int64 `json:"winner"`
		Loser  int64 `json:"loser"`
	}

	results := Results{a.state.GetWinner(), a.state.GetLoser()}
	message, _ := json.Marshal(results)

	for i := range a.clients {
		msg := NewGameMessage(messageEndGame, string(message))
		a.clients[i].conn.respond(msg)
		a.clients[i].notifications <- eventEndGame
	}
}

func (a *Arena) start(match *model.Game, moveChan chan Move) {
	a.match = match
	a.state.StartWithDifficulty(10)

	// TODO: DO better than this, this will race
	for i := range a.clients {
		a.clients[i].matchID = match.ID
		a.clients[i].moveChan = moveChan
	}

	a.notifyClients(eventStartGame)
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
	// Negative positions are no-ops moves
	if nextPos.X > -0 {
		a.state.MoveUser(playerID, nextPos)
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
