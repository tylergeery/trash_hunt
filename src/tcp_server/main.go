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

	client := connection.NewClient(conn)
	err := client.ValidateUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	// TODO: this should be an event to manager
	manager.AddPending(client)

	// TODO: let player know when arena is set
	client.WaitForStart()
}
