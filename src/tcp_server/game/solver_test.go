package game

import (
	"testing"
)

func TestEasySolverCanSolve(t *testing.T) {
	// Given
	p1, p2 := NewPlayer(1), NewPlayer(2)
	p1.Solver = NewSolver(1)
	state := NewState(p1, p2)
	state.StartWithDifficulty(1)

	totalPossibleMoves := (8 * 8 * 4) // squares with 4 DOF (middle board)
	totalPossibleMoves += (8 * 4 * 3) // positions with 3  (side walls)
	totalPossibleMoves += (4 * 2)     // positions only have 2 DOF (corners)
	solved := false

	for i := 0; i < totalPossibleMoves; i++ {
		nextPos := p1.Solver.GetMove(1, state)
		state.MoveUser(1, nextPos)
		if state.GetWinner() == 1 {
			solved = true
			break
		}
	}

	if !solved {
		t.Fatalf("Game was not solved by EasySolver")
	}
}
