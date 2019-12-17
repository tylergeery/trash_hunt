package game

import (
	"testing"
	"time"
)

func TestEasySolverCanSolve(t *testing.T) {
	// Given
	ogTime := time.Unix(time.Now().Unix()-1, 0)
	p1, p2 := NewPlayer(1), NewPlayer(2)
	p1.Solver = NewSolver(1)
	state := NewState(p1, p2)
	state.StartWithDifficulty(1)

	totalPossibleMoves := (8 * 8 * 4) // squares with 4 DOF (middle board)
	totalPossibleMoves += (8 * 4 * 3) // positions with 3  (side walls)
	totalPossibleMoves += (4 * 2)     // positions only have 2 DOF (corners)
	totalPossibleMoves *= 4           // possibly would need to backtrack in each direction
	solved := false

	for i := 0; i < totalPossibleMoves; i++ {
		nextPos := p1.Solver.GetMove(1, state)
		p1.lastMoveTime = ogTime
		state.MoveUser(p1.ID, nextPos)

		if state.GetWinner() == 1 {
			solved = true
			break
		}
	}

	if !solved {
		t.Fatalf("Game was not solved by EasySolver")
	}

	if state.GetLoser() != 2 {
		t.Fatalf("Expected Player 2 to be loser")
	}
}
