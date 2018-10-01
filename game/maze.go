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
	m.removeWalls(m.recent)
}

func (m *Maze) addWalls(pos int) {
	m.recent = pos
	m.setWalls(pos, true)
}

func (m *Maze) removeWalls(pos int) {
	m.setWalls(pos, false)
}

func (m *Maze) setWalls(pos int, value bool) {
	wall := (pos % 4)

	switch wall {
	case 0: // top wall
		m.setWall(pos, value)
		if (pos / gameBoardSize) != 0 {
			m.setWall(pos-(gameBoardSize*4)+2, value)
		}
	case 1: // right wall
		m.setWall(pos, value)
		if (pos % gameBoardSize) < ((gameBoardSize - 1) * 4) {
			m.setWall(pos+6, value)
		}
	case 2: // bottom wall
		m.setWall(pos, value)
		if (pos / gameBoardSize) != (gameBoardSize - 1) {
			m.setWall(pos+(gameBoardSize*4)-2, value)
		}
	case 3: // left wall
		m.setWall(pos, value)
		if (pos % gameBoardSize) > 4 {
			m.setWall(pos-6, value)
		}
	}
}

func (m *Maze) setWall(pos int, value bool) {
	rowTotal := 4 * gameBoardSize

	m.Walls[pos/rowTotal][(pos%rowTotal)/4][pos%4] = value
}
