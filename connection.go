package main

import (
	"net"
	"time"
)

type Connection struct {
	conn net.Conn
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) Read() ([]byte, error) {
	buf := make([]byte, 1024)

	n, errRead := c.conn.Read(buf)
	if errRead != nil {
		return nil, errRead
	}

	return buf[:n], nil
}

func (c *Connection) Write(data []byte) error {
	_, err := c.conn.Write(data)

	return err
}

func (c *Connection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Connection) Close() {
	_ = c.conn.Close()
}
