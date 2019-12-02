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
	// set up game manager
	manager = connection.NewManager(connectionCount)

	// start manager routine for ng matching clients
	go manager.Start()

	// TODO: make sure the manager will clean up lingering goroutines

	// handle websocket connections
	go acceptSocketConnections(manager)

	// main routine
	acceptTCPConnections(manager)
}

func acceptSocketConnections(manager *connection.Manager) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Socket request received")
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("Socket upgrade error: %s", err.Error())
			return
		}

		connOpen := true
		conn.SetCloseHandler(func(code int, text string) error {
			connOpen = false
			return nil
		})
		conn.SetReadDeadline(time.Now().Add(10 * time.Minute)) // TODO: game duration
		defer conn.Close()

		manager.InitCh <- connection.NewSocketConnection(conn)

		for connOpen {
			// TODO: Do better, maybe a lock or waitGroup
			time.Sleep(10 * time.Millisecond)
		}
		fmt.Println("Letting go of websocket connection")
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}

func acceptTCPConnections(manager *connection.Manager) {
	service := ":8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		panic(err)
	}

	// set up thread pool
	for i := 0; i < connectionCount; i++ {
		go handleConnection(manager)
	}

	// accept connections
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}

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
