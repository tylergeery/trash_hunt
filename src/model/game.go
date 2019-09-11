package model

import (
	"github.com/tylergeery/trash_hunt/storage"
)

const GameStatusPending = 1
const GameStatusActive = 2
const GameStatusComplete = 3

var gameCreateColumns = []string{
	"player1_id",
	"player2_id",
	"status",
}
var gameQueryColumns = []string{
	"id",
	"player1_id",
	"player2_id",
	"winner_id",
	"loser_id",
	"status",
	"ended_at",
	"created_at",
	"updated_at",
}
var gameUpdateColumns = []string{
	"winner_id",
	"loser_id",
	"status",
	"ended_at",
}

type Game struct {
	ID        int64  `json:"id"`
	Player1ID int64  `json:"player1_id"`
	Player2ID int64  `json:"player2_id"`
	WinnerID  int64  `json:"winner_id"`
	LoserID   int64  `json:"loser_id"`
	Status    int    `json:"status"`
	EndedAt   string `json:"ended_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// NewGame returns a new game object
func NewGame(player1_id, player2_id int64) *Game {
	return &Game{
		Player1ID: player1_id,
		Player2ID: player2_id,
		Status:    GameStatusActive,
	}
}

// CreateNewGame creates and saves a new game objet
func CreateNewGame(player1_id, player2_id int64) (*Game, error) {
	game := NewGame(player1_id, player2_id)

	err := game.Save()

	return game, err
}

func (g *Game) Save() error {
	var id int64
	var err error

	if g.ID == 0 {
		id, err = storage.Insert(storage.TABLE_GAME, g.toCreateValues(), gameCreateColumns)
		g.ID = id
	} else {
		err = storage.Update(storage.TABLE_GAME, g.toUpdateValues(), gameUpdateColumns, g.ID)
	}

	return err
}

func (g *Game) toCreateValues() []interface{} {
	return []interface{}{g.Player1ID, g.Player2ID, g.Status}
}

func (g *Game) toUpdateValues() []interface{} {
	return []interface{}{g.WinnerID, g.LoserID, g.Status, g.EndedAt}
}
