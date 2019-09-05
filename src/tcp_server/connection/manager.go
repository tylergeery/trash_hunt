package connection

import (
	"fmt"
	"net"

	"github.com/tylergeery/trash_hunt/model"
)

// Manager handles the creation of matches (updating movements)
type Manager struct {
	ActiveCh  chan *Arena
	ExitCh    chan *Client
	InitCh    chan *net.TCPConn
	PendingCh chan *Client
	active    map[int64]*Arena
	pending   map[int64]*Client
}

// NewManager returns a new connection manager object
func NewManager(clients int) *Manager {
	return &Manager{
		ActiveCh:  make(chan *Arena, clients),
		ExitCh:    make(chan *Client, clients),
		InitCh:    make(chan *net.TCPConn, clients),
		PendingCh: make(chan *Client, clients),
		active:    map[int64]*Arena{},
		pending:   map[int64]*Client{},
	}
}

// Start initializes a manager waiting for events
func (m *Manager) Start() {
	m.waitForEvents()
}

func (m *Manager) waitForEvents() {
	for {
		select {
		case client := <-m.PendingCh:
			m.addPending(client)
			fmt.Printf("Manager: adding client %s\n", client)
			m.match()
		case client := <-m.ExitCh:
			m.removePending(client)
			fmt.Printf("Manager: removing client: %s\n", client)
			m.endMatch(client.matchID) // TODO: move this as activeCh move
		case _ar := <-m.ActiveCh:
			// TODO: Update arena with move
			fmt.Printf("Manager: got move from client: %s\n", _ar)
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
	fmt.Println("Manager: trying to match")
	ids := []int64{}

	for k := range m.pending {
		ids = append(ids, k)
	}

	// TODO: make a real matching algorithm
	for len(ids) > 1 {
		client1, _ := m.pending[ids[0]]
		client2, _ := m.pending[ids[1]]
		arena := NewArena(client1, client2)
		game, err := model.CreateNewGame(client1.player.ID, client2.player.ID)
		if err != nil {
			fmt.Printf("Manager: could not start match (%s)\n", err)
			continue
		}

		fmt.Printf("Manager: created game (%s)\n", game.ID)

		arena.notifyClients(1)

		fmt.Printf("Manager: done notifying clients\n")

		// update manager's records
		m.active[game.ID] = arena
		m.removePending(client1).removePending(client2)
		ids = ids[2:]
	}

	return m
}
