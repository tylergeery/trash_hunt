package game

import (
	"testing"
)

func TestNewMaze(t *testing.T) {
	maze := NewMaze()

	for i := 0; i < gameBoardSize; i++ {
		for j := 0; j < gameBoardSize; j++ {
			for k := 0; k < 4; k++ {
				exp := int8(0)
				if i == 0 && k == 0 {
					exp = 1
				} else if i == (gameBoardSize-1) && k == 2 {
					exp = 1
				} else if j == 0 && k == 3 {
					exp = 1
				} else if j == (gameBoardSize-1) && k == 1 {
					exp = 1
				}

				if maze.Walls[i][j][k] != exp {
					t.Fatalf("Wall %d %d %d is not %d", i, j, k, exp)
				}
			}
		}
	}
}

func TestAddWallsAndRevert(t *testing.T) {
	// make test cases
	type TestCase struct {
		positions []int
		walls     [][][3]int
	}

	cases := []TestCase{
		TestCase{
			positions: []int{0, 1, 2, 3},
			walls: [][][3]int{
				[][3]int{
					[3]int{0, 0, 0},
				},
				[][3]int{
					[3]int{0, 0, 1},
					[3]int{0, 1, 3},
				},
				[][3]int{},
				[][3]int{},
			},
		},
		TestCase{
			positions: []int{57, 55, 56, 54},
			walls: [][][3]int{
				[][3]int{},
				[][3]int{},
				[][3]int{},
				[][3]int{},
			},
		},
		TestCase{
			positions: []int{98, 96, 99, 97},
			walls: [][][3]int{
				[][3]int{},
				[][3]int{},
				[][3]int{},
				[][3]int{},
			},
		},
	}

	// setup test vars
	maze := NewMaze()

	for _, c := range cases {
		for i, pos := range c.positions {
			maze.addWalls(pos)
			for _, wall := range c.walls[i] {
				// check that walls exist
				if maze.Walls[wall[0]][wall[1]][wall[2]] == 0 {
					t.Fatalf("Wall should exist at spot maze.Walls[%d][%d][%d]", wall[0], wall[1], wall[2])
				}
			}

			maze.revert()
			for _, wall := range c.walls[i] {
				// check that walls no longer exist
				if maze.Walls[wall[0]][wall[1]][wall[2]] == 1 {
					t.Fatalf("Wall should not exist at spot maze.Walls[%d][%d][%d]: %d", wall[0], wall[1], wall[2], maze.Walls[wall[0]][wall[1]][wall[2]])
				}
			}
		}
	}
}

func TestCanMoveBetween(t *testing.T) {
	type TestCase struct {
		pos1 Pos
		pos2 Pos
		exp  bool
	}

	edge := gameBoardSize - 1
	output := map[bool]string{true: "to be able", false: "to not be able"}
	maze := NewMaze()
	testCases := []TestCase{
		TestCase{
			pos1: Pos{0, 0},
			pos2: Pos{0, 0},
			exp:  false,
		},
		TestCase{
			pos1: Pos{1, 3},
			pos2: Pos{1, 5},
			exp:  false,
		},
		TestCase{
			pos1: Pos{3, 5},
			pos2: Pos{5, 5},
			exp:  false,
		},
		TestCase{
			pos1: Pos{edge - 1, edge - 1},
			pos2: Pos{edge, edge},
			exp:  false,
		},
		TestCase{
			pos1: Pos{edge, edge},
			pos2: Pos{edge, edge + 1},
			exp:  false,
		},
		TestCase{
			pos1: Pos{1, 1},
			pos2: Pos{1, 2},
			exp:  true,
		},
		TestCase{
			pos1: Pos{5, 6},
			pos2: Pos{6, 6},
			exp:  true,
		},
		TestCase{
			pos1: Pos{4, edge},
			pos2: Pos{4, edge - 1},
			exp:  true,
		},
		TestCase{
			pos1: Pos{1, 6},
			pos2: Pos{2, 6},
			exp:  true,
		},
		TestCase{
			pos1: Pos{edge, 1},
			pos2: Pos{edge, 2},
			exp:  true,
		},
	}

	for _, test := range testCases {
		if ok := maze.CanMoveBetween(test.pos1, test.pos2); ok != test.exp {
			t.Fatalf("Expected %s to move between (%d, %d) and (%d, %d)", output[test.exp], test.pos1.X, test.pos1.Y, test.pos2.X, test.pos2.Y)
		}
	}
}
