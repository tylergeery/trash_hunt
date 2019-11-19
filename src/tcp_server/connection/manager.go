package connection

import (
	"fmt"

	"github.com/tylergeery/trash_hunt/model"
	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

// Manager handles the creation of matches (updating movements)
type Manager struct {
	ActiveCh  chan Move
	ExitCh    chan *Client
	InitCh    chan Connection
	PendingCh chan *Client
	active    map[int64]*Arena
	pending   map[int64]*Client
}

// NewManager returns a new connection manager object
func NewManager(clients int) *Manager {
	return &Manager{
		ActiveCh:  make(chan Move, clients),
		ExitCh:    make(chan *Client, clients),
		InitCh:    make(chan Connection, clients),
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
			fmt.Printf("Manager: adding client: (%d)\n", client.player.ID)
			m.addPending(client).match()
		case client := <-m.ExitCh:
			fmt.Printf("Manager: removing client:(%d)\n", client.player.ID)
			m.removePending(client).endMatch(client.matchID)
		case move := <-m.ActiveCh:
			fmt.Printf("Manager: got move from client: (%d) (%d, %d)\n", move.playerID, move.pos.X, move.pos.Y)
			arena, ok := m.active[move.matchID]
			if !ok {
				fmt.Printf("Manager: got move for unknown game (%d)", move.matchID)
				continue
			}

			arena.moveUser(move.playerID, move.pos)
			arena.sendPositions()
		}
	}
}

func (m *Manager) addPending(client *Client) *Manager {
	// Check if user wants to play computer
	if client.preferences.Opponent == -1 {
		m.createMatch(client, nil, 1)
	}

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

func (m *Manager) createMatch(client1, client2 *Client, difficulty int) bool {
	var arena *Arena
	var match *model.Game
	var err error

	if client2 == nil {
		// Create Computer
		computer := game.NewPlayer(-1)
		computer.Solver = game.NewSolver(difficulty)

		arena = NewArena(client1.player, computer, client1)
		match, err = model.CreateNewGame(client1.player.ID, -1)
	} else {
		arena = NewArena(client1.player, client2.player, client1, client2)
		match, err = model.CreateNewGame(client1.player.ID, client2.player.ID)
	}

	if err != nil {
		fmt.Printf("Manager: could not start match (%s)\n", err)
		return false
	}

	arena.state.StartWithDifficulty(10)
	fmt.Printf("Manager: created game (%d)\n", match.ID)

	arena.start(match.ID, m.ActiveCh)

	fmt.Printf("Manager: done notifying clients\n")

	// update manager's records
	m.active[match.ID] = arena

	return true
}

// Match takes pending players and creates matches
func (m *Manager) match() *Manager {
	playersByDifficulty := map[string][]int64{}

	for k := range m.pending {

		diff := m.pending[k].preferences.Difficulty
		playersByDifficulty[diff] = append(playersByDifficulty[diff], k)
	}

	// TODO: make a real matching algorithm
	for d := range playersByDifficulty {
		for len(playersByDifficulty[d]) > 1 {
			client1, _ := m.pending[playersByDifficulty[d][0]]
			client2, _ := m.pending[playersByDifficulty[d][1]]
			success := m.createMatch(client1, client2, 1)
			if !success {
				continue
			}
			m.removePending(client1).removePending(client2)
			playersByDifficulty[d] = playersByDifficulty[d][2:]
		}
	}

	return m
}
