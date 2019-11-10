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
	var gameMessage connection.GameMessage
	message := connection.ReadStringFromConn(conn, make([]byte, 50))
	_ = json.Unmarshal([]byte(message), &gameMessage)

	if strings.TrimPrefix(gameMessage.Data, "Status: ") != "Pending" {
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
	var gameMessage connection.GameMessage
	_ = json.Unmarshal([]byte(resp), &gameMessage)

	var state, state2 game.State
	_ = json.Unmarshal([]byte(gameMessage.Data), &state)

	moves := []byte("lrud")
	playerMoved := false
	pos := state.Players[player.ID].GetPos()
	for _, move := range moves {
		fmt.Println("move", string(move))
		_, err := playerConn.Write([]byte{move})
		if err != nil {
			panic(fmt.Sprintf("Error sending left move: %s", err))
		}

		nextPos := game.Pos{pos.X, pos.Y}
		switch move {
		case 'l':
			nextPos.X--
		case 'r':
			nextPos.X++
		case 'u':
			nextPos.Y++
		case 'd':
			nextPos.Y--
		}

		positions := connection.ReadStringFromConn(playerConn, make([]byte, 200))
		_ = json.Unmarshal([]byte(positions), &gameMessage)
		_ = json.Unmarshal([]byte(gameMessage.Data), &state2.Players)

		if state2.Players[player.ID].Pos.X != nextPos.X || state2.Players[player.ID].Pos.Y != nextPos.Y {
			fmt.Printf(
				"Player could not move (%s) from Pos(%d, %d) to Pos(%d, %d): %s",
				move, pos.X, pos.Y, nextPos.X, nextPos.Y, state.Maze,
			)
			continue
		}

		playerMoved = true
		break
	}

	if !playerMoved {
		panic(
			fmt.Sprintf(
				"Player was expected to move from pos (%d, %d) : %s",
				state.Players[player.ID].GetPos().X, state.Players[player.ID].GetPos().Y,
				state.Maze,
			),
		)
	}

	//  Make sure the computer player moved as well
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
