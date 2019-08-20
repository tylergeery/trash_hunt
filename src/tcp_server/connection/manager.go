package connection

import (
	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

// Manager handles connected players in games
type Manager struct {
	games   map[int64]*game.State
	pending map[int64]*Client
}

// NewManager returns a new connection manager object
func NewManager() *Manager {
	return &Manager{
		games:   make(map[int64]*game.State),
		pending: make(map[int64]*Client),
	}
}

// AddPending adds a new player to the pending queue
func (m *Manager) AddPending(client *Client) *Manager {
	m.pending[client.player.ID] = client

	return m
}

func (c *Manager) removePending(id int64) *Manager {
	delete(c.pending, id)

	return c
}

func (c *Manager) endMatch(id int64) *Manager {
	_, _ = c.games[id]

	// cleanup

	delete(c.games, id)

	return c
}

// Match takes pending players and creates games
func (c *Manager) Match() *Manager {
	// TODO: find appropriate matches
	// move matches from pending to active

	return c
}
