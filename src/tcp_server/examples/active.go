package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/model"
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
	conn.SetReadDeadline(time.Now().Add(10 * time.Minute)) // TODO: game duration
	defer conn.Close()

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

	// Wait for game to match
	player1Response := make([]byte, 50)
	player2Response := make([]byte, 50)
	player2Conn.Read(player2Response)
	player1Conn.Read(player1Response)
	game1 := string(bytes.TrimRight(player1Response, "\x00"))
	game2 := string(bytes.TrimRight(player2Response, "\x00"))
	if game1 != game2 {
		panic(fmt.Sprintf("Players were not paired for the same game: Player1 %s, Player2 %s", game1, game2))
	}

	// Move one character and ensure the other receives the move

	// Move them both and ensure each receives the other

	// Try to win?

}
