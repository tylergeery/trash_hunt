package model

import (
	"testing"
)

func TestCreateAndUpdateGame(t *testing.T) {
	game, err := CreateNewGame(1, 2)
	if err != nil {
		t.Fatalf("Unexpected error creating game: %s", err.Error())
	}

	if game.ID == 0 {
		t.Fatalf("Game ID should not be 0")
	}

	gameFromID, err := GameFromID(game.ID)
	if err != nil {
		t.Fatalf("Unexpected error calling GameFromID: %s", err.Error())
	}

	if game.ID != gameFromID.ID {
		t.Fatalf("Did not receive match game ID. Expected (%d), got (%d)", game.ID, gameFromID.ID)
	}
}
