package game

import (
	"fmt"
	"sort"
)

// Solver provides a common interface for all automated game solvers
type Solver interface {
	GetMove(playerID int64, state *State) Pos
}

// EasySolver solves the game in the simplest (slow) way
type EasySolver struct {
	tracked map[string]int
}

// GetMove returns a preferred move for the EasySolver
func (s *EasySolver) GetMove(playerID int64, state *State) Pos {
	availableMoves := state.GetAvailableMoves(playerID)

	sort.Slice(availableMoves, func(i, j int) bool {
		return s.moveCount(availableMoves[i]) < s.moveCount(availableMoves[j])
	})

	return s.recordMove(availableMoves[0])
}

func (s *EasySolver) moveCount(pos Pos) int {
	key := fmt.Sprintf("%d,%d", pos.X, pos.Y)

	return s.tracked[key]
}

func (s *EasySolver) recordMove(pos Pos) Pos {
	key := fmt.Sprintf("%d,%d", pos.X, pos.Y)
	s.tracked[key]++

	return pos
}

// NewSolver returns a Solver interface based on difficulty
func NewSolver(difficulty int) Solver {
	return &EasySolver{
		tracked: map[string]int{},
	}
}
