package game

import (
	"errors"

	"github.com/tylergeery/trash_hunt/storage"
)

const dbTable = "player"

// Player is a given player in the game
type Player struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	pw         string
	Name       string `json:"name"`
	FacebookID string `json:"facebook_id"`
	Pos        Pos    `json:"pos"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

var types = map[string]string{
	"id":          "int",
	"email":       "string",
	"name":        "string",
	"facebook_id": "string",
	"created_at":  "string",
	"updated_at":  "string",
}

// NewPlayer - Constructor
func PlayerNew(id int64, email, pw, name, facebookID, createdAt, updatedAt string) *Player {
	return &Player{
		ID:         id,
		Email:      email,
		pw:         pw,
		Name:       name,
		FacebookID: facebookID,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

// Register a new player
func PlayerRegister(email, pw, name, facebookID string) (*Player, error) {
	// Validate and hash password, or validate facebookID

	// Validate email is unique

	p := PlayerNew(0, email, pw, name, facebookID, "", "")
	err := p.save()

	return p, err
}

// Update an existing player
func (p *Player) PlayerUpdate() error {
	if p.ID == 0 {
		return errors.New("Could not update non-existent player")
	}

	return p.save()
}

func (p *Player) save() error {
	var err error
	var id int64

	err = p.validate()

	if err != nil {
		return err
	}

	if p.ID == 0 {
		id, err = storage.Insert(dbTable, p.toSaveMap(), types)
		p.ID = id
	} else {
		err = storage.Update(dbTable, p.toSaveMap(), types)
	}

	return err
}

func (p *Player) toSaveMap() map[string]interface{} {
	return map[string]interface{}{
		"id":    p.ID,
		"name":  p.Name,
		"email": p.Email,
	}
}

func (p *Player) validate() error {
	var err error

	// check valid email

	// return error
	return err
}
