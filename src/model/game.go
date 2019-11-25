package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tylergeery/trash_hunt/storage"
)

const (
	// GameStatusPending is pending game status
	GameStatusPending = 1
	// GameStatusActive is active game status
	GameStatusActive = 2
	// GameStatusComplete is completed game status
	GameStatusComplete = 3
)

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

// Game model record also referred to as "match"
type Game struct {
	ID        int64     `json:"id"`
	Player1ID int64     `json:"player1_id"`
	Player2ID int64     `json:"player2_id"`
	WinnerID  int64     `json:"winner_id"`
	LoserID   int64     `json:"loser_id"`
	Status    int       `json:"status"`
	EndedAt   time.Time `json:"ended_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewGame returns a new game object
func NewGame(player1ID, player2ID int64) *Game {
	return &Game{
		Player1ID: player1ID,
		Player2ID: player2ID,
		Status:    GameStatusActive,
	}
}

// CreateNewGame creates and saves a new game objet
func CreateNewGame(player1ID, player2ID int64) (*Game, error) {
	game := NewGame(player1ID, player2ID)

	err := game.Save()

	return game, err
}

// GameFromID gets a Game Model from an ID
func GameFromID(id int64) (*Game, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(gameQueryColumns, ","),
		storage.TABLE_GAME,
	)

	game := scanGame(storage.FetchRow(query, id))

	if game.ID == 0 {
		return nil, fmt.Errorf("Could not find game with id: %d", id)
	}

	return game, nil
}

// Save updates existing game record, or creates a new one
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

func scanGame(row *sql.Row) *Game {
	var id, player1ID, player2ID, winnerID, loserID int64
	var status int
	var endedAt, createdAt, updatedAt time.Time

	row.Scan(&id, &player1ID, &player2ID, &winnerID, &loserID, &status, &endedAt, &createdAt, &updatedAt)

	return &Game{
		ID:        id,
		Player1ID: player1ID,
		Player2ID: player2ID,
		WinnerID:  winnerID,
		LoserID:   loserID,
		Status:    status,
		EndedAt:   endedAt,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
