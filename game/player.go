package game

// Player is a given player in the game
type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Pos  Pos    `json:"pos"`
}
