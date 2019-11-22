package connection

import "testing"

func TestManagerCreatesMatches(t *testing.T) {
	clientCount := 4
	manager := NewManager(clientCount)
	difficulties := []string{"", "easy", "medium", "difficult"}

	go manager.Start()

	// add some clients, but no preferences should match
	for i := 1; i < clientCount; i++ {
		client := NewClient(nil)
		client.preferences = GameSetUp{
			Difficulty: difficulties[i],
		}
		manager.PendingCh <- client
	}

	if len(manager.active) != 0 {
		t.Fatalf("Manager created unexpected active matches: %d", len(manager.active))
	}

	client := NewClient(nil)
	client.preferences = GameSetUp{
		Difficulty: "easy",
	}
	manager.PendingCh <- client

	// if len(manager.active) != 1 {
	// 	t.Fatalf("Manager did not create expected easy match: %d, pending: %d", len(manager.active), len(manager.pending))
	// }
}

func TestManagerCreateComputerMatch(t *testing.T) {
	// TODO
}
