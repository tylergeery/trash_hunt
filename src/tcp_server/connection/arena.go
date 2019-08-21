package connection

import (
    "github.com/tylergeery/trash_hunt/tcp_server/game"
)

// Arena holds all the information about clients and the active game
type Arena struct {
    state *game.State
    clients []*Client
}

// NewArena sets up all the infrastructure for gameplay
func NewArena(c1, c2 *Client) *Arena {
    // Set up channels to communicate to both clients

    // Set up channels to receive events from both clients

    return &Arena{
        state: game.NewState(c1.player, c2.player),
        clients: []*Client{c1, c2},
    }
}
