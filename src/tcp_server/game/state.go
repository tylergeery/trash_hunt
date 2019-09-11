package game

import (
	"fmt"
	"math/rand"
	"time"
)

const gameBoardSize = 10

// State controls all state related to the game playing
type State struct {
	Players map[int64]*Player `json:"players"`
	Maze    *Maze             `json:"maze"`
}

// NewState sets new maze and player/trash pos for gameplay
func NewState(player1, player2 *Player) *State {
	var s State

	// Initialize a new random seed
	rand.Seed(time.Now().UnixNano())

	// initialize random trash position
	s.Maze = NewMaze()
	s.Maze.TrashPos.X = 1
	s.Maze.TrashPos.Y = 9
	s.Players = make(map[int64]*Player)
	s.Players[player1.ID] = player1
	player1.Pos.X = rand.Intn(gameBoardSize)
	player1.Pos.Y = rand.Intn(gameBoardSize)
	s.Players[player2.ID] = player2
	player2.Pos.X = rand.Intn(gameBoardSize)
	player2.Pos.Y = rand.Intn(gameBoardSize)

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
	// Players can share outcomes and visited state
	outcomes := map[string]bool{}
	visited := []Pos{}

	for id := range s.Players {
		player := s.Players[id]
		if !s.PlayerCanFinish(player, outcomes, visited) {
			return false
		}
	}

	return true
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

// TODO: test
// MoveUser changes the current position of a user to the nextPos
func (s *State) MoveUser(playerID int64, nextPos Pos) {
	player := s.Players[playerID]

	if s.Maze.CanMoveBetween(player.Pos, nextPos) {
		player.Pos = nextPos
	}
}

func hasVisited(pos Pos, visited []Pos) bool {
	for _, v := range visited {
		if v.X == pos.X && v.Y == pos.Y {
			return true
		}
	}

	return false
}
