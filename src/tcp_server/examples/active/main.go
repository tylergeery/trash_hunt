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

	playerPreferences := connection.GameSetUp{
		UserToken:  playerToken,
		Opponent:   0,      // 0 means no preference
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
	fmt.Println("Starting game between two active players")

	// Create two players
	player1 := model.GetTestPlayer("active1")
	player2 := model.GetTestPlayer("active2")

	// Setup two players
	player1Conn := setupPlayer(player1)
	player2Conn := setupPlayer(player2)
	defer player1Conn.Close()
	defer player2Conn.Close()

	// Wait for game to match
	str1 := connection.ReadStringFromConn(player1Conn, make([]byte, 2500))
	str2 := connection.ReadStringFromConn(player2Conn, make([]byte, 2500))

	var gameMessage1, gameMessage2 connection.GameMessage
	_ = json.Unmarshal([]byte(str1), &gameMessage1)
	_ = json.Unmarshal([]byte(str2), &gameMessage2)

	if gameMessage1.Data != gameMessage2.Data {
		panic(fmt.Sprintf("Players were not paired for the same game: Player1 %s, Player2 %s", gameMessage1.Data, gameMessage2.Data))
	}

	var state1, state2 game.State
	_ = json.Unmarshal([]byte(gameMessage1.Data), &state1)
	_ = json.Unmarshal([]byte(gameMessage2.Data), &state2)
	fmt.Println(str1, gameMessage1.Data)

	// Move one character and ensure the other receives the move
	moves := []byte("lrud")
	playerMoved := false
	pos := state2.Players[player1.ID].GetPos()
	for _, move := range moves {
		_, err := player1Conn.Write([]byte{move})
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
			nextPos.Y--
		case 'd':
			nextPos.Y++
		}

		positions1 := connection.ReadStringFromConn(player1Conn, make([]byte, 200))
		positions2 := connection.ReadStringFromConn(player2Conn, make([]byte, 200))
		_ = json.Unmarshal([]byte(positions1), &gameMessage1)
		_ = json.Unmarshal([]byte(positions2), &gameMessage2)

		_ = json.Unmarshal([]byte(gameMessage1.Data), &state1.Players)
		_ = json.Unmarshal([]byte(gameMessage2.Data), &state2.Players)
		if gameMessage1.Data != gameMessage2.Data {
			panic(
				fmt.Sprintf(
					"Players did not received the same game updates: Player1 %s, Player2 %s",
					gameMessage1.Data, gameMessage2.Data,
				),
			)
		}

		if state1.Players[player1.ID].Pos.X != nextPos.X || state1.Players[player1.ID].Pos.Y != nextPos.Y {
			fmt.Printf(
				"Player could not move (%s) from Pos(%d, %d) to Pos(%d, %d): %+v",
				string(move), pos.X, pos.Y, nextPos.X, nextPos.Y, state1.Maze,
			)
			continue
		}

		playerMoved = true
		break
	}

	if !playerMoved {
		panic(
			fmt.Sprintf(
				"Player1 was expected to move from pos (%d, %d):\n %+v",
				state1.Players[player1.ID].GetPos().X, state1.Players[player1.ID].GetPos().Y,
				state2.Maze,
			),
		)
	}

	// TODO: Let one player leave and make sure results are recorded correctly

	// Try to win?
	fmt.Println("Game over")
}
