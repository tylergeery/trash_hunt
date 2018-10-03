package game

import (
	"testing"
)

func TestInitializeGameState(t *testing.T) {
	p1 := PlayerNew(1, "", "", "", "", "", "")
	p2 := PlayerNew(1, "", "", "", "", "", "")
	state := InitializeGameState(p1, p2)

	if state.Maze.TrashPos.X == 0 && state.Maze.TrashPos.Y == 0 {
		t.Fatalf("Trash Pos still at (0,0)")
	}

	if state.Player1.Pos.X == 0 && state.Player1.Pos.Y == 0 {
		t.Fatalf("Player1 Pos still at (0,0)")
	}

	if state.Player2.Pos.X == 0 && state.Player2.Pos.Y == 0 {
		t.Fatalf("Player2 Pos still at (0,0)")
	}
}

func TestGetAvailableMoves(t *testing.T) {
	type TestCase struct {
		p   Pos
		exp []Pos
	}

	p1 := PlayerNew(1, "", "", "", "", "", "")
	p2 := PlayerNew(1, "", "", "", "", "", "")
	state := InitializeGameState(p1, p2)
	testCases := []TestCase{
		TestCase{
			p:   Pos{0, 0},
			exp: []Pos{Pos{1, 0}, Pos{0, 1}},
		},
	}

	for _, test := range testCases {
		p1.Pos = test.p
		moves := state.getAvailableMoves(p1)

		if len(moves) != len(test.exp) {
			t.Fatalf("Moves had len: %d, but expected length: %d", len(moves), len(test.exp))
		}

		for i := range moves {
			if moves[i].X != test.exp[i].X {
				t.Fatalf("Expected move %d to have Pos (%d, %d), but had (%d, %d)", i, test.exp[i].X, test.exp[i].Y, moves[i].X, moves[i].Y)
			}
		}
	}
}
