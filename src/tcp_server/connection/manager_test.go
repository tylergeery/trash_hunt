package connection

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

func TestManagerCreatesMatches(t *testing.T) {
	clientCount := 4
	manager := NewManager(clientCount)
	matchedClients := []*Client{}
	difficulties := []string{"", "easy", "medium", "difficult"}
	var wg sync.WaitGroup
	wg.Add(2)

	go manager.Start()

	// add some clients, but no preferences should match
	for i := 1; i < clientCount; i++ {
		client := NewClient(nil)
		client.player = game.NewPlayer(int64(i))
		client.conn = &MockConnection{}
		client.preferences = GameSetUp{
			Difficulty: difficulties[i],
		}
		manager.PendingCh <- client

		if i == 1 {
			matchedClients = append(matchedClients, client)
		}
	}

	if len(manager.active) != 0 {
		t.Fatalf("Manager created unexpected active matches: %d", len(manager.active))
	}

	client := NewClient(nil)
	client.player = game.NewPlayer(5)
	client.conn = &MockConnection{}
	client.preferences = GameSetUp{
		Difficulty: "easy",
	}
	manager.PendingCh <- client
	matchedClients = append(matchedClients, client)

	for i := range matchedClients {
		go func(client *Client) {
			select {
			case event := <-client.notifications:
				fmt.Println("Client notification:", event)
				wg.Done()
			}
		}(matchedClients[i])
	}

	wg.Wait()

	time.Sleep(1 * time.Second)
	if len(manager.active) != 1 {
		t.Fatalf("Manager did not create expected easy match: %d, pending: %d", len(manager.active), len(manager.pending))
	}
}

func TestManagerCreateComputerMatch(t *testing.T) {
	// TODO
}
