package game

// Pos handles a position on the game board
type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Maze handles all the information related to the game maze
type Maze struct {
	TrashPos Pos `json:"trashPos"`
}

// State controls all state related to the game playing
type State struct {
	Player1 Player `json:"player1"`
	Player2 Player `json:"player2"`
	Maze    Maze   `json:"maze"`
}
