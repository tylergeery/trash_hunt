package game

import (
	"fmt"
	"math/rand"
	"time"
)

const gameBoardSize = 15

// State controls all state related to the game playing
type State struct {
	Players map[int64]*Player `json:"players"`
	Maze    *Maze             `json:"maze"`
	done    bool
	winner  int64
	loser   int64
}

// NewState sets new maze and player/trash pos for gameplay
func NewState(player1, player2 *Player) *State {
	var s State

	// Initialize a new random seed
	rand.Seed(time.Now().UnixNano())

	// initialize random trash position
	s.Maze = NewMaze()
	s.Maze.TrashPos.X = rand.Intn(gameBoardSize)
	s.Maze.TrashPos.Y = rand.Intn(gameBoardSize)

	s.Players = make(map[int64]*Player)
	s.Players[player1.GetID()] = player1
	s.Players[player2.GetID()] = player2

	// TODO: Do better to try and make the games "fair"
	for true {
		pos1 := Pos{rand.Intn(gameBoardSize), rand.Intn(gameBoardSize)}
		pos2 := Pos{rand.Intn(gameBoardSize), rand.Intn(gameBoardSize)}

		if pos1.X == pos2.X && pos1.Y == pos2.Y {
			continue
		}

		if pos1.X == s.Maze.TrashPos.X && pos1.Y == s.Maze.TrashPos.Y {
			continue
		}

		if pos2.X == s.Maze.TrashPos.X && pos2.Y == s.Maze.TrashPos.Y {
			continue
		}

		player1.setPos(pos1)
		player2.setPos(pos2)
		break
	}

	return &s
}

// StartWithDifficulty creates new game state with specified difficulty
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

// IsValid checks if maze is valid
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

// PlayerCanFinish checks if the given player finish
func (s *State) PlayerCanFinish(player *Player, outcomes map[string]bool, visited []Pos) bool {
	if player.GetPos().X == s.Maze.TrashPos.X && player.GetPos().Y == s.Maze.TrashPos.Y {
		return true
	}

	originalPos := player.GetPos()
	defer func() {
		player.setPos(originalPos)
	}()

	for _, pos := range s.findAvailableMoves(player) {
		if hasVisited(pos, visited) {
			continue
		}

		player.setPos(pos)
		key := fmt.Sprintf("%d-%d", pos.X, pos.Y)
		visited = append(visited, pos)

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

// GetAvailableMoves return positions available for given playerID
func (s *State) GetAvailableMoves(playerID int64) []Pos {
	return s.findAvailableMoves(s.Players[playerID])
}

func (s *State) findAvailableMoves(player *Player) []Pos {
	positions := []Pos{}
	possibles := []Pos{
		Pos{player.GetPos().X, player.GetPos().Y - 1}, // up
		Pos{player.GetPos().X + 1, player.GetPos().Y}, // right
		Pos{player.GetPos().X, player.GetPos().Y + 1}, // down
		Pos{player.GetPos().X - 1, player.GetPos().Y}, // left
	}

	for _, next := range possibles {
		if next.Y < 0 || next.Y >= gameBoardSize {
			continue
		}

		if next.X < 0 || next.X >= gameBoardSize {
			continue
		}

		// check for blocking walls
		if !s.Maze.CanMoveBetween(player.GetPos(), next) {
			continue
		}

		positions = append(positions, next)
	}

	return positions
}

// MoveUser changes the current position of a user to the nextPos
func (s *State) MoveUser(playerID int64, nextPos Pos) error {
	if s.done {
		return fmt.Errorf("game is already over")
	}

	player := s.Players[playerID]

	// Ensure that they don't move more than 1 time per sec
	if time.Now().Sub(player.lastMoveTime).Nanoseconds() < 100000000 {
		fmt.Printf("State: player (%d) moving too fast: %s, %s\n", playerID, player.lastMoveTime, time.Now())
		return fmt.Errorf("player moving too quickly")
	}

	if !s.Maze.CanMoveBetween(player.GetPos(), nextPos) {
		return nil
	}

	player.setPos(nextPos)
	player.lastMoveTime = time.Now()
	if s.hasWon(player.ID) {
		s.SetWinner(player.ID)
	}

	return nil
}

// SetWinner sets the winner of game
func (s *State) SetWinner(winner int64) {
	s.winner = winner
	s.done = true

	for id := range s.Players {
		if winner != id {
			s.loser = id
			break
		}
	}
}

// SetLoser sets the loser of game
func (s *State) SetLoser(loser int64) {
	s.loser = loser
	s.done = true

	for id := range s.Players {
		if loser != id {
			s.winner = id
			break
		}
	}
}

// GetWinner of game
func (s *State) GetWinner() int64 {
	if !s.done {
		return 0
	}

	return s.winner
}

// GetLoser of game
func (s *State) GetLoser() int64 {
	if !s.done {
		return 0
	}

	return s.loser
}

func (s *State) hasWon(playerID int64) bool {
	if s.done {
		return s.winner == playerID
	}

	pos := s.Players[playerID].Pos
	if pos.X == s.Maze.TrashPos.X && pos.Y == s.Maze.TrashPos.Y {
		return true
	}

	return false
}

func hasVisited(pos Pos, visited []Pos) bool {
	for _, v := range visited {
		if v.X == pos.X && v.Y == pos.Y {
			return true
		}
	}

	return false
}
