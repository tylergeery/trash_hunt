package game

// Player holds the temporary game-related player info
type Player struct {
	ID          int64 `json:"id"`
	Pos         Pos   `json:"pos"`
	preferences map[string]string
}

// NewPlayer creates a game Player object
func NewPlayer(id int64) *Player {
	return &Player{
		ID:          id,
		Pos:         Pos{},
		preferences: make(map[string]string),
	}
}
