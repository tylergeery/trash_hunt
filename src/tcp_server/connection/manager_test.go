package connection

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/tylergeery/trash_hunt/model"
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

func TestManagerEndsMatch(t *testing.T) {
	// Given
	manager := NewManager(5)
	player1, player2 := game.NewPlayer(200), game.NewPlayer(300)
	client := NewClient(&MockConnection{})
	client.player = player1
	arena := NewArena(player1, player2, client)
	match, _ := model.CreateNewGame(player1.ID, player2.ID)
	fmt.Println("Created match:", match.ID)
	arena.match = match
	manager.active[match.ID] = arena
	client.matchID = match.ID
	endTime := time.Now()

	// When
	go manager.Start()
	manager.ExitCh <- client

	// Then
	time.Sleep(2 * time.Second)
	game, err := model.GameFromID(match.ID)
	if err != nil {
		t.Fatalf("Unexpected game fetch error: %s", err.Error())
	}

	if game.WinnerID != player2.ID {
		t.Fatalf("Game winner (%d) was not expected (%d)", game.WinnerID, player2.ID)
	}

	if game.LoserID != player1.ID {
		t.Fatalf("Game loser (%d) was not expected (%d)", game.LoserID, player1.ID)
	}

	if game.Status != model.GameStatusComplete {
		t.Fatalf("Game status (%d) was not expected (%d)", game.Status, model.GameStatusComplete)
	}

	if game.EndedAt.Unix() >= endTime.Unix() {
		t.Fatalf("Game time (%s) was not >= expected (%s)", game.EndedAt, endTime)
	}
}
