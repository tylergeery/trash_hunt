package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/tylergeery/trash_hunt/tcp_server/connection"
)

const connectionCount = 100

var manager *connection.Manager

func main() {
	service := ":8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		panic(err)
	}

	// set up game manager
	manager = connection.NewManager(connectionCount)

	// set up thread pool
	for i := 0; i < connectionCount; i++ {
		go handleConnection(manager)
	}

	// accept connections
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}

	// start manager routine for ng matching clients
	manager.Start()

	// pass connections to pool
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Connection accept error: %s", err.Error())
			continue
		}
		manager.InitCh <- conn
	}
}

func handleConnection(manager *connection.Manager) {
	for {
		conn := <-manager.InitCh

		fmt.Printf("Got a connection: %s\n", conn)
		handleClient(conn, manager)
	}
}

func handleClient(conn net.Conn, manager *connection.Manager) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Minute)) // TODO: game duration
	defer conn.Close()

	// Try to turn connection into game player
	client := connection.NewClient(conn)
	err := client.ValidateUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	// Make sure we let the manager know the client is gone
	defer func() {
		manager.ExitCh <- client
	}()

	// Client will wait for manager to start a match
	manager.PendingCh <- client
	client.WaitForStart()
}
