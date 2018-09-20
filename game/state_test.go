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

func TestInitializeMaze(t *testing.T) {
	maze := NewMaze()

	maze.Initialize()

	if maze.TrashPos.X == 0 && maze.TrashPos.Y == 0 {
		t.Fatalf("Trash Pos still at (0,0) %s", maze.TrashPos)
	}
}
