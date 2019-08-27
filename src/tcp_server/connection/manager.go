package connection

import (
	"fmt"
	"net"
)

// Manager handles the creation of matches (updating movements)
type Manager struct {
	ActiveCh  chan *Arena
	ExitCh    chan *Client
	InitCh    chan net.Conn
	PendingCh chan *Client
	active    map[int64]*Arena
	pending   map[int64]*Client
}

// NewManager returns a new connection manager object
func NewManager(clients int) *Manager {
	return &Manager{
		ActiveCh:  make(chan *Arena, clients),
		ExitCh:    make(chan *Client, clients),
		InitCh:    make(chan net.Conn, clients),
		PendingCh: make(chan *Client, clients),
		active:    map[int64]*Arena{},
		pending:   map[int64]*Client{},
	}
}

// Start initializes a manager waiting for events
func (m *Manager) Start() {
	// TODO: wait for new connections
	go m.waitForEvents()

	// TODO; Update computer players, make matches
}

func (m *Manager) waitForEvents() {
	for {
		select {
		case client := <-m.PendingCh:
			m.addPending(client)
			fmt.Printf("Manager: adding client %s", client)
		case client := <-m.PendingCh:
			m.removePending(client)
			fmt.Printf("Manager: removing client: %s", client)
			m.endMatch(client.matchID) // TODO: decide who lost
		case _ = <-m.ActiveCh:
			// TODO: Update arena with move
		}
	}
}

func (m *Manager) addPending(client *Client) *Manager {
	m.pending[client.player.ID] = client

	return m
}

func (m *Manager) removePending(client *Client) *Manager {
	delete(m.pending, client.player.ID)

	return m
}

func (m *Manager) endMatch(matchID int64) *Manager {
	_, _ = m.active[matchID]

	// cleanup

	delete(m.active, matchID)

	return m
}

// Match takes pending players and creates matches
func (m *Manager) match() *Manager {
	// TODO: find appropriate matches
	// move matches from pending to active

	return m
}
