package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/model"
	"github.com/tylergeery/trash_hunt/tcp_server/connection"
	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

func setupPlayer(player *model.Player) net.Conn {
	// Create tokens for player
	playerToken, err := auth.CreateToken(player)
	if err != nil {
		panic(fmt.Sprintf("Could not create token for player: %s", err))
	}

	// set up player preferences
	playerPreferences := connection.GameSetUp{
		UserToken:  playerToken,
		Opponent:   -1,     // -1 means computer opponent
		Difficulty: "easy", // To help with matching
	}
	preferences, err := json.Marshal(playerPreferences)
	if err != nil {
		panic(fmt.Sprintf("Could not create player preferences: %s", err))
	}

	// open game connection
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(fmt.Sprintf("Could not open connection: %s", err))
	}

	// setting connection settings
	conn.SetReadDeadline(time.Now().Add(10 * time.Second)) // TODO: game duration

	// authenticate
	fmt.Printf("Sending player preferences: %s\n", playerToken)
	_, err = conn.Write(preferences)
	if err != nil {
		panic(fmt.Sprintf("Could not write player preferences: %s", playerToken))
	}

	// listen for reply
	message := connection.ReadStringFromConn(conn, make([]byte, 50))

	if strings.TrimPrefix(message, "Status: ") != "Pending" {
		panic(fmt.Sprintf("Received something other than status pending: \n%s\n%s", string(message), strings.TrimPrefix(string(message), "Status: ")))
	}

	fmt.Printf("Player established: %s\n", playerToken)
	return conn
}

func main() {
	fmt.Println("Starting game against computer.")

	// Create two players
	player := model.GetTestPlayer("active3")

	// Setup two players
	playerConn := setupPlayer(player)
	defer playerConn.Close()

	// Wait for game to match
	resp := connection.ReadStringFromConn(playerConn, make([]byte, 1500))
	var state game.State
	json.Unmarshal([]byte(resp), &state)

	// Move one character and ensure the other receives the move
	move := 'l'
	state1Player1 := state.Players[player.ID]
	fmt.Println(state.Players)
	if state1Player1.GetPos().X == 0 {
		move = 'r'
	}
	_, err := playerConn.Write([]byte{byte(move)})
	if err != nil {
		panic(fmt.Sprintf("Error sending left move: %s", err))
	}
	positions := connection.ReadStringFromConn(playerConn, make([]byte, 100))
	state2 := game.State{
		Players: make(map[int64]*game.Player),
	}
	fmt.Println(positions)
	json.Unmarshal([]byte(positions), &state2.Players)

	pos := state2.Players[player.ID].GetPos().X + 1
	if move == 'r' {
		pos = state2.Players[player.ID].GetPos().X - 1
	}
	if pos != state.Players[player.ID].GetPos().X {
		panic(
			fmt.Sprintf(
				"Player was expected to move (%s) from pos (%d, %d) to pos (%d, %d)",
				string(move),
				state.Players[player.ID].GetPos().X, state.Players[player.ID].GetPos().Y,
				state2.Players[player.ID].GetPos().X, state2.Players[player.ID].GetPos().Y,
			),
		)
	}

	// TODO: Move sure the computer player moved as well
	if state.Players[-1].GetPos().X == state2.Players[-1].GetPos().X && state.Players[-1].GetPos().Y == state2.Players[-1].GetPos().Y {
		panic(
			fmt.Sprintf(
				"Computer Player was expected to move from pos (%d, %d)",
				state.Players[-1].GetPos().X, state.Players[-1].GetPos().Y,
			),
		)
	}

	// Try to win?
	fmt.Println("Game over")
}
