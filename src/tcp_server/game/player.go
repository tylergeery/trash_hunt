package game

// Player holds the user game-related player info
type Player struct {
	ID          int64 `json:"id"`
	Pos         Pos   `json:"pos"`
	preferences map[string]string
	Solver      Solver `json:"-"`
}

// NewPlayer creates a game Player object
func NewPlayer(id int64) *Player {
	return &Player{
		ID:          id,
		Pos:         Pos{},
		preferences: make(map[string]string),
		Solver:      nil,
	}
}

// GetID return player ID
func (p *Player) GetID() int64 {
	return p.ID
}

// GetPos return player Pos
func (p *Player) GetPos() Pos {
	return p.Pos
}

func (p *Player) setPos(pos Pos) {
	p.Pos = pos
}
