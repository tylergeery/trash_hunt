package game

// Maze handles all the information related to the game maze
type Maze struct {
	TrashPos Pos        `json:"trashPos"`
	Walls    [][][]bool `json:"walls"`
	recent   int
}

// NewMaze - maze constructor
func NewMaze() *Maze {
	walls := make([][][]bool, gameBoardSize)

	for i := 0; i < gameBoardSize; i++ {
		walls[i] = make([][]bool, gameBoardSize)

		for j := 0; j < gameBoardSize; j++ {
			walls[i][j] = []bool{false, false, false, false}
		}
	}

	return &Maze{
		TrashPos: Pos{},
		Walls:    walls,
	}
}

func (m *Maze) revert() {

}

func (m *Maze) addWalls(pos int) {
	wall := (pos % 4)

	switch wall {
	case 0: // top wall
		m.addWall(pos)
		if (pos / gameBoardSize) != 0 {
			m.addWall(pos - (gameBoardSize * 4) + 2)
		}
	case 1: // right wall
		m.addWall(pos)
		if (pos % gameBoardSize) < ((gameBoardSize - 1) * 4) {
			m.addWall(pos + 6)
		}
	case 2: // bottom wall
		m.addWall(pos)
		if (pos / gameBoardSize) != (gameBoardSize - 1) {
			m.addWall(pos + (gameBoardSize * 4) - 2)
		}
	case 3: // left wall
		m.addWall(pos)
		if (pos % gameBoardSize) > 4 {
			m.addWall(pos - 6)
		}
	}
}

func (m *Maze) addWall(pos int) {
	rowTotal := 4 * gameBoardSize

	m.Walls[pos/rowTotal][(pos%rowTotal)/4][pos%4] = true
}

// IsValid - ensure maze is valid
func (m *Maze) IsValid() bool {
	// Check that player1 can reach trash

	// Check that player2 can reach trash

	return true
}
