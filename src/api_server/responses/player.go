package responses

import (
	"github.com/tylergeery/trash_hunt/src/game"
)

// PlayerUpdateResponse - Response to PlayerUpdateRequest
type PlayerUpdateResponse struct {
	Token string `json:"token"`
}

// PlayerLoginResponse - Response to PlayerLoginRequest
type PlayerLoginResponse struct {
	Player *game.Player `json:"player"`
	Token  string       `json:"token"`
}
