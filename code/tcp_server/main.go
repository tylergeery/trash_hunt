package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/tylergeery/trash_hunt/game"
)

const connectionCount = 10

func main() {
	service := ":8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	// set up thread pool
	events := make(chan net.Conn, connectionCount*2)
	for i := 0; i < connectionCount; i++ {
		go handleConnection(events)
	}

	// accept connetions
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// pass connections to pool
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Connection accept error: %s", err.Error())
			continue
		}
		events <- conn
	}
}

func handleConnection(connections chan net.Conn) {
	for {
		conn := <-connections

		// set up what game we are in
		game := game.State{}

		handleClient(conn, game)
	}
}

func handleClient(conn net.Conn, gameState game.State) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Minute)) // set 2 minutes timeout
	request := make([]byte, 10)
	defer conn.Close()

	for {
		readLen, err := conn.Read(request)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Connection read error: %s", err.Error())
			break
		}

		if readLen == 0 {
			break // connection already closed by client
		}

		switch request[0] {
		case 'l':
			// move left
		case 'r':
			// move right
		case 'u':
			// move up
		case 'd':
			// move down
		}

		// update for opponent status

		// check for winner

		msg, err := json.Marshal(gameState)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Marshal error: %s", err.Error())
			continue
		}
		conn.Write(msg)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
