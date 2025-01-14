package server

import (
	"net"
	"time"
)

type connection struct {
	conn net.Conn
}

func newConnection(conn net.Conn) *connection {
	return &connection{
		conn: conn,
	}
}

func (c *connection) Read() ([]byte, error) {
	buf := make([]byte, 1024)

	n, errRead := c.conn.Read(buf)
	if errRead != nil {
		return nil, errRead
	}

	return buf[:n], nil
}

func (c *connection) Write(data []byte) (int, error) {
	return c.conn.Write(data)
}

func (c *connection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *connection) Close() {
	_ = c.conn.Close()
}
