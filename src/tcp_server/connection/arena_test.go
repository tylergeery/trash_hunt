package connection

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/tylergeery/trash_hunt/model"
	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

// MockConnection object used for mocking Connection interface while testing arena
type MockConnection struct {
	input   string
	message GameMessage
}

func (c *MockConnection) gatherInput(input []byte) (s string, err error) {
	return c.input, nil
}

func (c *MockConnection) respond(message GameMessage) error {
	c.message = message

	return nil
}

func TestStartGame(t *testing.T) {
	p1 := game.NewPlayer(15)
	p2 := game.NewPlayer(150)
	clients := []*Client{}
	clientCount := 3
	notificationChannel := make(chan int, 3)

	for i := 0; i < clientCount; i++ {
		client := NewClient(&MockConnection{})
		clients = append(clients, client)
		go func() {
			select {
			case move := <-client.notifications:
				notificationChannel <- move
			}
		}()
	}

	match := &model.Game{
		ID: 100,
	}
	arena := NewArena(p1, p2, clients...)
	arena.start(match, make(chan Move, 10))

	for i := range clients {
		conn := clients[i].conn.(*MockConnection)
		if conn.message.Event != messageInitGame {
			t.Fatalf("Client received unexpected game message: %d", conn.message.Event)
		}
	}

	time.Sleep(2 * time.Second)
	if len(notificationChannel) != 3 {
		t.Fatalf("Expected all clients to notified of start game")
	}

	close(notificationChannel)
	for event := range notificationChannel {
		if event != eventStartGame {
			t.Fatalf("Client got unexpected notification: %d", event)
		}
	}
}

func TestCompleteGame(t *testing.T) {
	// Given
	p1 := game.NewPlayer(47)
	p2 := game.NewPlayer(1147)
	c1 := NewClient(&MockConnection{})
	c1.player = p1
	c2 := NewClient(&MockConnection{})
	c2.player = p2
	state := game.State{
		Players: make(map[int64]*game.Player),
		Maze:    game.NewMaze(),
	}
	state.Players[c1.player.GetID()] = c1.player
	state.Players[c2.player.GetID()] = c2.player
	state.Maze.TrashPos = game.Pos{1, 1}
	p1.Pos = game.Pos{1, 0}
	match := model.NewGame(p1.ID, p2.ID)
	arena := Arena{&state, match, []*Client{c1, c2}}

	// When first user wins
	arena.moveUser(p1.ID, game.Pos{1, 1})
	arena.end()

	// Then
	c1Message, c2Message := c1.conn.(*MockConnection).message, c2.conn.(*MockConnection).message
	if c1Message.Event != messageEndGame || c2Message.Event != messageEndGame {
		t.Fatalf("Expected end game message: user 1 got %d, user 2 got %d", c1Message.Event, c2Message.Event)
	}

	if c1Message.Data != c2Message.Data {
		t.Fatalf("Expected clients to receive same message: %s, %s", c1Message.Data, c2Message.Data)
	}

	var results map[string]int64
	_ = json.Unmarshal([]byte(c1Message.Data), &results)

	if results["winner"] != p1.ID || results["loser"] != p2.ID {
		t.Fatalf("Unexpected results: %s", c1Message.Data)
	}
}
