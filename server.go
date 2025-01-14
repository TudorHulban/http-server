package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	bufferPool sync.Pool
}

func NewServer() *Server {
	return &Server{
		bufferPool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(nil)
			},
		},
	}
}

func (s *Server) Run(address string) error {
	listener, errListen := net.Listen("tcp", address)
	if errListen != nil {
		return fmt.Errorf(
			"listener start: %w",
			errListen,
		)
	}
	defer listener.Close()

	log.Printf("Listening on %s...", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)

			continue
		}

		go s.handleConnection(
			NewConnection(conn),
		)
	}
}

func (s *Server) handleConnection(conn *Connection) {
	defer conn.Close()

	// Set a read deadline for idle connections
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	for {
		if errOnTraffic := s.onTraffic(conn); errOnTraffic != nil {

			break
		}

		// Reset the timeout after successful activity
		conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	}
}

func (s *Server) SendStatus(statusCode int) []byte {
	buffer := s.bufferPool.Get().(*bytes.Buffer)
	defer s.bufferPool.Put(buffer)
	buffer.Reset()

	buffer.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode)))
	buffer.WriteString("Content-Length: 0\r\n") // No body, so Content-Length is 0
	buffer.WriteString("\r\n")                  // Blank line to terminate headers

	return buffer.Bytes()
}

func (s *Server) SendBody(statusCode int, body string) []byte {
	buffer := s.bufferPool.Get().(*bytes.Buffer)
	defer s.bufferPool.Put(buffer)
	buffer.Reset()

	buffer.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode)))
	buffer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
	buffer.WriteString("\r\n") // Blank line between headers and body
	buffer.WriteString(body)

	return buffer.Bytes()
}

func (s *Server) onTraffic(conn *Connection) error {
	data, errRead := conn.Read()
	if errRead != nil {
		return errRead
	}

	_, errParse := NewHTTPRequest(string(data))
	if errParse != nil {
		log.Printf("failed to parse HTTP request: %v", errParse)

		_ = conn.Write(
			[]byte(
				"HTTP/1.1 400 Bad Request\r\nContent-Length: 0\r\n\r\n",
			),
		)

		return nil
	}

	// go log.Printf("Received HTTP request: %v", request)

	return conn.Write(
		s.SendStatus(
			http.StatusOK,
		),
	)
}
