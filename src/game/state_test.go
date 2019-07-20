package game

import (
	"testing"
)

func TestInitializeGameState(t *testing.T) {
	p1 := PlayerNew(1, "", "", "", "", "", PlayerStatusActive, "", "")
	p2 := PlayerNew(1, "", "", "", "", "", PlayerStatusActive, "", "")
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

func TestInitializeWithDifficulty(t *testing.T) {
	difficulties := []int{1, 5, 10}

	for _, d := range difficulties {
		p1 := PlayerNew(1, "", "", "", "", "", PlayerStatusActive, "", "")
		p2 := PlayerNew(1, "", "", "", "", "", PlayerStatusActive, "", "")
		state := InitializeGameState(p1, p2)
		state.StartWithDifficulty(d)

		if !state.IsValid() {
			t.Fatalf("Received invalid game state")
		}
	}
}

func TestPlayerCanFinish(t *testing.T) {
	p1 := PlayerNew(1, "", "", "", "", "", PlayerStatusActive, "", "")
	p2 := PlayerNew(1, "", "", "", "", "", PlayerStatusActive, "", "")
	state := InitializeGameState(p1, p2)
	state.Player1.Pos.X = state.Maze.TrashPos.X
	state.Player1.Pos.Y = state.Maze.TrashPos.Y
	state.Player2.Pos.X = state.Maze.TrashPos.X
	state.Player2.Pos.Y = state.Maze.TrashPos.Y - 1
	outcomes := map[string]bool{}
	visited := []Pos{}

	if !state.PlayerCanFinish(p1, outcomes, visited) {
		t.Fatalf("Player at pos (%d, %d) could not finish", p1.Pos.X, p1.Pos.Y)
	}

	if !state.PlayerCanFinish(p2, outcomes, visited) {
		t.Fatalf("Player at pos (%d, %d) could not finish", p2.Pos.X, p2.Pos.Y)
	}
}

func TestGetAvailableMoves(t *testing.T) {
	type TestCase struct {
		p   Pos
		exp []Pos
	}

	p1 := PlayerNew(1, "", "", "", "", "", PlayerStatusActive, "", "")
	p2 := PlayerNew(1, "", "", "", "", "", PlayerStatusActive, "", "")
	state := InitializeGameState(p1, p2)
	testCases := []TestCase{
		TestCase{
			p:   Pos{0, 0},
			exp: []Pos{Pos{1, 0}, Pos{0, 1}},
		},
		TestCase{
			p:   Pos{1, 0},
			exp: []Pos{Pos{2, 0}, Pos{1, 1}, Pos{0, 0}},
		},
		TestCase{
			p:   Pos{9, 1},
			exp: []Pos{Pos{9, 0}, Pos{9, 2}, Pos{8, 1}},
		},
		TestCase{
			p:   Pos{0, 9},
			exp: []Pos{Pos{0, 8}, Pos{1, 9}},
		},
		TestCase{
			p:   Pos{5, 5},
			exp: []Pos{Pos{5, 4}, Pos{6, 5}, Pos{5, 6}, Pos{4, 5}},
		},
	}

	for _, test := range testCases {
		p1.Pos = test.p
		visited := []Pos{}
		moves := state.getAvailableMoves(p1, visited)

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
