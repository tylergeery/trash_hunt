package game

import (
	"github.com/tylergeery/trash_hunt/storage"
)

const dbTable = "player"

// Player is a given player in the game
type Player struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Pos   Pos    `json:"pos"`
}

var types = map[string]string{
	"id":    "int",
	"email": "string",
	"name":  "string",
	"pos":   "string",
}

// Save - Update or Create player
func (p *Player) Save() error {
	var err error

	if p.ID == 0 {
		var id int64
		id, err = storage.Insert(dbTable, p.toSaveMap(), types)
		if err != nil {
			return err
		}
		p.ID = id
	} else {
		storage.Update(dbTable, p.toSaveMap(), types)
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
