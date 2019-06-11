package game

import (
	"fmt"
	"math/rand"
	"time"
)

const gameBoardSize = 10

// State controls all state related to the game playing
type State struct {
	Player1 *Player `json:"player1"`
	Player2 *Player `json:"player2"`
	Maze    *Maze   `json:"maze"`
}

// InitializeGameState - setup new maze and player/trash pos for gameplay
func InitializeGameState(player1, player2 *Player) *State {
	var s State

	// Initialize a new random seed
	rand.Seed(time.Now().UnixNano())

	// initialize random trash position
	s.Maze = NewMaze()
	s.Maze.TrashPos.X = 1
	s.Maze.TrashPos.Y = 9
	s.Player1 = player1
	s.Player1.Pos.X = rand.Intn(gameBoardSize)
	s.Player1.Pos.Y = rand.Intn(gameBoardSize)
	s.Player2 = player2
	s.Player2.Pos.X = rand.Intn(gameBoardSize)
	s.Player2.Pos.Y = rand.Intn(gameBoardSize)

	return &s
}

// StartWithDifficulty - Creates new game state with specified difficulty
func (s *State) StartWithDifficulty(difficulty int) {
	// keep adding walls as long as both user can solve
	for i := 0; i < difficulty; {
		// try adding a new wall
		newWall := rand.Intn(gameBoardSize*(gameBoardSize-1)*4 + gameBoardSize*4)
		s.Maze.addWalls(newWall)

		// Check maze is still solvable for both players
		if !s.IsValid() {
			i++

			// Undo latest wall and try again
			s.Maze.revert()
		}
	}
}

// IsValid - ensure maze is valid
func (s *State) IsValid() bool {
	outcomes := map[string]bool{}
	visited := []Pos{}

	return s.PlayerCanFinish(s.Player1, outcomes, visited) && s.PlayerCanFinish(s.Player2, outcomes, visited)
}

// PlayerCanFinish - can the given player finish?
func (s *State) PlayerCanFinish(player *Player, outcomes map[string]bool, visited []Pos) bool {
	// fmt.Printf("Player can finish (%d, %d), trash (%d, %d)\n", player.Pos.X, player.Pos.Y, s.Maze.TrashPos.X, s.Maze.TrashPos.Y)
	if player.Pos.X == s.Maze.TrashPos.X && player.Pos.Y == s.Maze.TrashPos.Y {
		return true
	}

	originalPosX := player.Pos.X
	originalPosY := player.Pos.Y
	defer func() {
		player.Pos.X = originalPosX
		player.Pos.Y = originalPosY
	}()

	for _, pos := range s.getAvailableMoves(player, visited) {
		player.Pos.X = pos.X
		player.Pos.Y = pos.Y
		key := fmt.Sprintf("%d-%d", pos.X, pos.Y)
		visited = append(visited, Pos{pos.X, pos.Y})

		if success, ok := outcomes[key]; ok {
			return success
		}

		success := s.PlayerCanFinish(player, outcomes, visited)

		// memoize results
		outcomes[key] = success

		if success {
			return true
		}
	}

	return false
}

func (s *State) getAvailableMoves(player *Player, visited []Pos) []Pos {
	positions := []Pos{}
	next := Pos{player.Pos.X, player.Pos.Y}

	// can player go up?
	next.Y = player.Pos.Y - 1
	if next.Y >= 0 && s.Maze.CanMoveBetween(player.Pos, next) && !hasVisited(next, visited) {
		positions = append(positions, Pos{next.X, next.Y})
	}
	next.Y = player.Pos.Y

	// can player go right?
	next.X = player.Pos.X + 1
	if next.X < gameBoardSize && s.Maze.CanMoveBetween(player.Pos, next) && !hasVisited(next, visited) {
		positions = append(positions, Pos{next.X, next.Y})
	}
	next.X = player.Pos.X

	// can player go down?
	next.Y = player.Pos.Y + 1
	if next.Y < gameBoardSize && s.Maze.CanMoveBetween(player.Pos, next) && !hasVisited(next, visited) {
		positions = append(positions, Pos{next.X, next.Y})
	}
	next.Y = player.Pos.Y

	// can player go left?
	next.X = player.Pos.X - 1
	if player.Pos.X >= 0 && s.Maze.CanMoveBetween(player.Pos, next) && !hasVisited(next, visited) {
		positions = append(positions, Pos{next.X, next.Y})
	}
	next.X = player.Pos.X

	return positions
}

func hasVisited(pos Pos, visited []Pos) bool {
	for _, v := range visited {
		if v.X == pos.X && v.Y == pos.Y {
			return true
		}
	}

	return false
}