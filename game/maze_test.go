package game

import (
	"testing"
)

func TestNewMaze(t *testing.T) {
	maze := NewMaze()

	for i := 0; i < gameBoardSize; i++ {
		for j := 0; j < gameBoardSize; j++ {
			for k := 0; k < 4; k++ {
				if maze.Walls[i][j][k] {
					t.Fatalf("Wall %d %d %d is not false", i, j, k)
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
				if !maze.Walls[wall[0]][wall[1]][wall[2]] {
					t.Fatalf("Wall should exist at spot %d %d %d", wall[0], wall[1], wall[2])
				}
			}

			maze.revert()
			for _, wall := range c.walls[i] {
				// check that walls no longer exist
				if maze.Walls[wall[0]][wall[1]][wall[2]] {
					t.Fatalf("Wall should not exist at spot %d %d %d", wall[0], wall[1], wall[2])
				}
			}
		}
	}
}
