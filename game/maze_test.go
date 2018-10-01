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

func TestRevert(t *testing.T) {

}

func TestAddWalls(t *testing.T) {

}

func TestIsValid(t *testing.T) {

}
