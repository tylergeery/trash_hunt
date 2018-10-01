package game

import (
	"math/rand"
	"time"
)

const gameBoardSize = 10
const difficulty = 10

// State controls all state related to the game playing
type State struct {
	Player1 *Player `json:"player1"`
	Player2 *Player `json:"player2"`
	Maze    *Maze   `json:"maze"`
}

// Initialize new maze for gameplay
func InitializeGameState() *State {
	var s State

	// Initialize a new random seed
	rand.Seed(time.Now().UnixNano())

	// initialize random trash position
	s.Maze = NewMaze()
	s.Maze.TrashPos.X = 1
	s.Maze.TrashPos.Y = 9

	// keep adding walls as long as both user can solve
	for i := 0; i < difficulty; {
		// try adding a new wall
		s.Maze.addWalls(rand.Intn(gameBoardSize * gameBoardSize))

		if !s.Maze.IsValid() {
			i++

			s.Maze.revert()
		}
	}

	return &s
}
