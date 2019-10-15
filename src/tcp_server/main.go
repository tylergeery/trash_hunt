package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
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
	go manager.Start()

	// TODO: make sure the manager will clean up lingering goroutines

	// handle websocket connections
	go acceptSocketConnections(manager)

	// pass connections to pool
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Connection accept error: %s", err.Error())
			continue
		}

		// handle cleanup
		conn.SetReadDeadline(time.Now().Add(10 * time.Minute)) // TODO: game duration
		conn.SetKeepAlive(true)
		defer conn.Close()

		manager.InitCh <- connection.NewTCPConnection(conn)
	}
}

func acceptSocketConnections(manager *connection.Manager) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Socket upgrade error: %s", err.Error())
			return
		}
		defer conn.Close()

		manager.InitCh <- connection.NewSocketConnection(conn)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}

func handleConnection(manager *connection.Manager) {
	for {
		conn := <-manager.InitCh

		fmt.Printf("Got a connection\n")
		handleClient(conn, manager)
	}
}

func handleClient(conn connection.Connection, manager *connection.Manager) {
	// Try to turn connection into game player
	client := connection.NewClient(conn)
	err := client.SetUpUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create TCP Client: %s", err.Error())
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
