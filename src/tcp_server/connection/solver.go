package connection

import (
	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

type Solver interface {
	GetMove(playerID int64, state *game.State) game.Pos
}

type EasySolver struct{}

func (s *EasySolver) GetMove(playerID int64, state *game.State) game.Pos {
	availableMoves := state.GetAvailableMoves(playerID)

	if len(availableMoves) > 0 {
		return availableMoves[0]
	}

	return state.Players[playerID].Pos
}

func NewSolver(difficulty int) Solver {
	return &EasySolver{}
}
