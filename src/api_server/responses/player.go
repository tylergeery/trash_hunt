package responses

import (
	"github.com/tylergeery/trash_hunt/src/game"
)

// PlayerCreateResponse - Response to PlayerCreateRequest
type PlayerCreateResponse struct {
	Token string `json:"token"`
}

// PlayerCreateResponse - Response to PlayerCreateRequest
type PlayerLoginResponse struct {
	Player *game.Player `json:"player"`
	Token  string       `json:"token"`
}
