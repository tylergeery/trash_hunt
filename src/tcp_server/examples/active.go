package main

import (
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

	if strings.TrimPrefix(string(message), "Status: ") != "Pending" {
		panic(fmt.Sprintf("Received something other than status pending: %s", string(message)))
	}

	return conn
}

func move() {

}

func main() {
	fmt.Println("Starting game between two active players")

	// Create two players
	player1 := model.PlayerNew(1, "", "", "", "", "", model.PlayerStatusActive, "", "")
	player2 := model.PlayerNew(2, "", "", "", "", "", model.PlayerStatusActive, "", "")

	// Setup two players
	player1Conn := setupPlayer(player1)
	player2Conn := setupPlayer(player2)
	fmt.Println(player1Conn, player2Conn)

	// Move one character and ensure the other receives the move

	// Move them both and ensure each receives the other

	// Try to win?

}
