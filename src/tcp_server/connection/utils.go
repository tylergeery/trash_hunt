package connection

import (
    "bytes"
    "fmt"
    "net"
)

func ReadStringFromConn(conn net.Conn, chars []byte) string {
    _, err :=conn.Read(chars)
    if err != nil {
        fmt.Printf("connection::ReadStringFromConn error: %s\n", err)
    }

    return string(bytes.TrimRight(chars, "\x00"))
}
