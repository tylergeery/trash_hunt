package connection

// Manager handles the creation of matches (updating movements)
type Manager struct {
	active  map[int64]*Arena
	pending map[int64]*Client
}

// NewManager returns a new connection manager object
func NewManager() *Manager {
	return &Manager{
		active:  make(map[int64]*Arena),
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
	_, _ = c.active[id]

	// cleanup

	delete(c.active, id)

	return c
}

// Match takes pending players and creates matches
func (c *Manager) Match() *Manager {
	// TODO: find appropriate matches
	// move matches from pending to active

	return c
}
