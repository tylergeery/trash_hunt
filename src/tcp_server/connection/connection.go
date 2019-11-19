package connection

import (
	"bytes"
	"fmt"
	"net"

	"github.com/gorilla/websocket"
)

// Connection is a generic connection interface for gameplay
type Connection interface {
	gatherInput(input []byte) (s string, err error)
	respond(message GameMessage) error
}

// TCPConnection object that implements the Connection interface
type TCPConnection struct {
	conn *net.TCPConn
}

// NewTCPConnection creates a new TCP connection
func NewTCPConnection(conn *net.TCPConn) *TCPConnection {
	return &TCPConnection{
		conn: conn,
	}
}

func (c *TCPConnection) gatherInput(input []byte) (s string, err error) {
	_, err = c.conn.Read(input)
	if err != nil {
		return
	}

	s = string(bytes.TrimRight(input, "\x00"))

	return
}

func (c *TCPConnection) respond(msg GameMessage) error {
	message := msg.ToBytes()
	_, err := c.conn.Write(message)
	if err != nil {
		fmt.Printf("TCPConnection: error sending response: %s\n", string(message))
	}

	return err
}

// SocketConnection object that implements the Connection interface
type SocketConnection struct {
	conn *websocket.Conn
}

// NewSocketConnection creates a new TCP connection
func NewSocketConnection(conn *websocket.Conn) *SocketConnection {
	return &SocketConnection{
		conn: conn,
	}
}

func (c *SocketConnection) gatherInput(input []byte) (s string, err error) {
	_, input, err = c.conn.ReadMessage()
	if err != nil {
		return
	}

	s = string(bytes.TrimRight(input, "\x00"))
	if s != string(input) {
		fmt.Println("SocketClient trimmed something")
	}

	return
}

func (c *SocketConnection) respond(msg GameMessage) error {
	message := msg.ToBytes()
	err := c.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		fmt.Printf("SocketClient: error sending response: %s\n", string(message))
	}

	return err
}
