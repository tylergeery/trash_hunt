package connection

import (
	"bytes"
	"fmt"
	"net"

	"github.com/gorilla/websocket"
)

type Connection interface {
	gatherInput(input []byte) (s string, err error)
	respond(message string) error
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

func (c *TCPConnection) respond(message string) error {
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Client: error sending response: %s\n", message)
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
	// TODO
	return
}

func (c *SocketConnection) respond(message string) error {
	// TODO
	return nil
}
