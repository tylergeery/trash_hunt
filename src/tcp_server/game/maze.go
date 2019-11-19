package game

// Maze handles all the information related to the game maze
type Maze struct {
	TrashPos Pos        `json:"trashPos"`
	Walls    [][][]int8 `json:"walls"`
	recent   int
}

func toInt8(b bool) int8 {
	if b {
		return 1
	}

	return 0
}

// NewMaze - maze constructor
func NewMaze() *Maze {
	walls := make([][][]int8, gameBoardSize)

	for i := 0; i < gameBoardSize; i++ {
		walls[i] = make([][]int8, gameBoardSize)

		for j := 0; j < gameBoardSize; j++ {
			walls[i][j] = []int8{
				toInt8(i == 0),
				toInt8(j == (gameBoardSize - 1)),
				toInt8(i == (gameBoardSize - 1)),
				toInt8(j == 0),
			}
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
	m.setWalls(pos, 1)
}

func (m *Maze) removeWalls(pos int) {
	m.setWalls(pos, 0)
}

func (m *Maze) setWalls(pos int, value int8) {
	wall := (pos % 4)

	switch wall {
	case 0: // top wall
		m.setWall(pos, value)
		if (pos / (gameBoardSize * 4)) != 0 {
			m.setWall(pos-(gameBoardSize*4)+2, value)
		}
	case 1: // right wall
		m.setWall(pos, value)
		if (pos % (gameBoardSize * 4)) < ((gameBoardSize - 1) * 4) {
			m.setWall(pos+6, value)
		}
	case 2: // bottom wall
		m.setWall(pos, value)
		if (pos / (gameBoardSize * 4)) != (gameBoardSize - 1) {
			m.setWall(pos+(gameBoardSize*4)-2, value)
		}
	case 3: // left wall
		m.setWall(pos, value)
		if (pos % (gameBoardSize * 4)) > 4 {
			m.setWall(pos-6, value)
		}
	}
}

func (m *Maze) setWall(pos int, value int8) {
	rowTotal := 4 * gameBoardSize

	m.Walls[pos/rowTotal][(pos%rowTotal)/4][pos%4] = value
}

func (m *Maze) hasWall(pos int) bool {
	rowTotal := 4 * gameBoardSize

	return (m.Walls[pos/rowTotal][(pos%rowTotal)/4][pos%4]) > 0
}

// CanMoveBetween two positons (are there walls blocking?)
func (m *Maze) CanMoveBetween(pos1, pos2 Pos) bool {
	if pos1.X != pos2.X && pos1.Y != pos2.Y {
		return false
	}

	if (pos1.X - pos2.X) == 1 {
		return !m.hasWall(gameBoardSize*4*pos1.Y + 4*pos1.X + 3)
	}

	if (pos1.X - pos2.X) == -1 {
		return !m.hasWall(gameBoardSize*4*pos1.Y + 4*pos1.X + 1)
	}

	if (pos1.Y - pos2.Y) == 1 {
		return !m.hasWall(gameBoardSize*4*pos1.Y + 4*pos1.X)
	}

	if (pos1.Y - pos2.Y) == -1 {
		return !m.hasWall(gameBoardSize*4*pos1.Y + 4*pos1.X + 2)
	}

	return false
}
