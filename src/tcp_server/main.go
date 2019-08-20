package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/tylergeery/trash_hunt/tcp_server/connection"
	"github.com/tylergeery/trash_hunt/tcp_server/game"
)

const connectionCount = 100

var manager *connection.Manager

func main() {
	service := ":8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		panic(err)
	}

	// set up thread pool
	events := make(chan net.Conn, connectionCount)
	for i := 0; i < connectionCount; i++ {
		go handleConnection(events)
	}

	// set up game manager
	manager = connection.NewManager()

	// accept connections
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}

	// start manager routine for ng matching clients
	go handleManager()

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

func handleManager() {
	for {
		manager.Match()

		time.Sleep(200 * time.Millisecond)
	}
}

func handleConnection(connections chan net.Conn) {
	for {
		conn := <-connections

		fmt.Println("Got a connection")
		handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Minute)) // TODO: game duration
	defer conn.Close()

	// TODO: get user ID
	// TODO: validate user ID
	player := game.NewPlayer(0)
	client := connection.NewClient(conn, player)
	manager.AddPending(client)

	request := make([]byte, 10)
	for {
		readLen, err := conn.Read(request)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Connection read error: %s", err.Error())
			break
		}

		if readLen == 0 {
			// TODO: award winner if necessary
			// TODO: close connection to opponent
			break // connection already closed by client
		}

		// validate movement has not been "too quick"

		switch request[0] {
		case 'l':
			// move left
		case 'r':
			// move right
		case 'u':
			// move up
		case 'd':
			// move down
		case 's':
			// start game
		case 'q':
			// quit, end game
		}

		// update for opponent status

		// check for winner

		// msg, err := json.Marshal(gameState)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Marshal error: %s", err.Error())
		// 	continue
		// }
		// conn.Write(msg)
		client.Respond()
	}
}
