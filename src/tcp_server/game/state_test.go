package game

import (
	"testing"
)

func TestNewState(t *testing.T) {
	p1 := NewPlayer(1)
	p2 := NewPlayer(1)
	state := NewState(p1, p2)

	if state.Maze.TrashPos.X == 0 && state.Maze.TrashPos.Y == 0 {
		t.Fatalf("Trash Pos still at (0,0)")
	}

	player1, _ := state.Players[p1.ID]
	player2, _ := state.Players[p2.ID]
	if player1.GetPos().X == 0 && player1.GetPos().Y == 0 {
		t.Fatalf("Player1 Pos still at (0,0)")
	}

	if player2.GetPos().X == 0 && player2.GetPos().Y == 0 {
		t.Fatalf("Player2 Pos still at (0,0)")
	}
}

func TestInitializeWithDifficulty(t *testing.T) {
	difficulties := []int{1, 5, 10}

	for _, d := range difficulties {
		p1 := NewPlayer(1)
		p2 := NewPlayer(1)
		state := NewState(p1, p2)
		state.StartWithDifficulty(d)

		if !state.IsValid() {
			t.Fatalf("Received invalid game state")
		}
	}
}

func TestPlayerCanFinish(t *testing.T) {
	p1 := NewPlayer(2)
	p2 := NewPlayer(10)
	state := NewState(p1, p2)
	p1.Pos.X = state.Maze.TrashPos.X
	p1.Pos.Y = state.Maze.TrashPos.Y
	p2.Pos.X = state.Maze.TrashPos.X
	p2.Pos.Y = state.Maze.TrashPos.Y - 1
	outcomes := map[string]bool{}
	visited := []Pos{}

	if !state.PlayerCanFinish(p1, outcomes, visited) {
		t.Fatalf("Player at pos (%d, %d) could not finish", p1.Pos.X, p1.Pos.Y)
	}

	if !state.PlayerCanFinish(p2, outcomes, visited) {
		t.Fatalf("Player at pos (%d, %d) could not finish", p2.Pos.X, p2.Pos.Y)
	}
}

func TestFindAvailableMoves(t *testing.T) {
	type TestCase struct {
		p   Pos
		exp []Pos
	}

	p1 := NewPlayer(1)
	p2 := NewPlayer(1)
	state := NewState(p1, p2)
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
		moves := state.findAvailableMoves(p1, visited)

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

func TestMoveUser(t *testing.T) {
	p1 := NewPlayer(3)
	p2 := NewPlayer(11)
	state := NewState(p1, p2)

	type TestCase struct {
		player   *Player
		orig     Pos
		next     Pos
		expected Pos
		result   bool
	}
	cases := []TestCase{
		TestCase{p1, Pos{0, 0}, Pos{1, 1}, Pos{0, 0}, false},
		TestCase{p1, Pos{0, 0}, Pos{0, 1}, Pos{0, 1}, true},
		TestCase{p2, Pos{9, 9}, Pos{9, 8}, Pos{9, 8}, true},
		TestCase{p2, Pos{9, 9}, Pos{3, 3}, Pos{9, 9}, false},
	}

	for _, c := range cases {
		c.player.setPos(c.orig)
		if state.MoveUser(c.player.ID, c.next) != c.result {
			t.Fatalf(
				"User moving from %+v to %+v, did not received expected result: %t",
				c.orig, c.next, c.result,
			)
		}

		if c.player.Pos.X != c.expected.X || c.player.Pos.Y != c.expected.Y {
			t.Fatalf(
				"User moving from %+v to %+v, did not end at expected %+v, but instead %+v",
				c.orig, c.next, c.expected, c.player.Pos,
			)
		}
	}
}
