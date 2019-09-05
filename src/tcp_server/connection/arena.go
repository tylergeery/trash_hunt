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
    return &Arena{
        state: game.NewState(c1.player, c2.player),
        clients: []*Client{c1, c2},
    }
}

func (a *Arena) notifyClients(move int) {
    for i := range a.clients {
        a.clients[i].notifications <- move
    }
}
