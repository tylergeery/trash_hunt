package main

import (
	"bytes"
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

	// open game connection
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(fmt.Sprintf("Could not open connection: %s", err))
	}

	// setting connection settings
	conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // TODO: game duration

	// authenticate
	fmt.Printf("Sending player token: %s\n", playerToken)
	_, err = conn.Write([]byte(playerToken))
	if err != nil {
		panic(fmt.Sprintf("Could not write player token: %s", playerToken))
	}

	// listen for reply
	message := make([]byte, 50)
	_, err = conn.Read(message)
	if err != nil {
		panic(fmt.Sprintf("Received error waiting for pending status: %s", err))
	}
	message = bytes.TrimRight(message, "\x00")

	if strings.TrimPrefix(string(message), "Status: ") != "Pending" {
		panic(fmt.Sprintf("Received something other than status pending: \n%s\n%s", string(message), strings.TrimPrefix(string(message), "Status: ")))
	}

	fmt.Printf("Player established: %s\n", playerToken)
	return conn
}

func move() {

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
	game1 := connection.ReadStringFromConn(player1Conn, make([]byte, 1500))
	game2 := connection.ReadStringFromConn(player2Conn, make([]byte, 1500))
	var state1 game.State
	json.Unmarshal([]byte(game1), &state1)
	if game1 != game2 {
		panic(fmt.Sprintf("Players were not paired for the same game: Player1 %s, Player2 %s", game1, game2))
	}

	// Move one character and ensure the other receives the move
	_, err := player1Conn.Write([]byte{'l'})
	if err != nil {
		panic(fmt.Sprintf("Error sending left move: %s", err))
	}
	game1 = connection.ReadStringFromConn(player1Conn, make([]byte, 1500))
	game2 = connection.ReadStringFromConn(player2Conn, make([]byte, 1500))
	var state2 game.State
	json.Unmarshal([]byte(game1), &state2)
	if game1 != game2 {
		panic(fmt.Sprintf("Players were not paired for the same game: Player1 %s, Player2 %s", game1, game2))
	}
	fmt.Println(state1, state2)
	if (state2.Player1.Pos.X + 1) != state1.Player1.Pos.X {
		panic(fmt.Sprintf("Player1 was expected to move left from pos (%s) to pos (%s)", state1.Player1.Pos, state2.Player1.Pos))
	}

	// Move them both and ensure each receives the other

	// Try to win?
	fmt.Println("Game over")
}
