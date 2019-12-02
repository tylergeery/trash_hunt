package connection

import (
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

func TestArena(t *testing.T) {
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
