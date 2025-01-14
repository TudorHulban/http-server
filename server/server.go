package server

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Server struct {
	bufferPool sync.Pool
	readerPool sync.Pool

	tlsConfig *tls.Config
}

func NewServer(certFile, keyFile string) (*Server, error) {
	cert, errLoadX509 := tls.LoadX509KeyPair(certFile, keyFile)
	if errLoadX509 != nil {
		return nil,
			fmt.Errorf("load x509 key pair: %w", errLoadX509)
	}

	return &Server{
			bufferPool: sync.Pool{
				New: func() interface{} {
					return bytes.NewBuffer(nil)
				},
			},
			readerPool: sync.Pool{
				New: func() interface{} {
					return bufio.NewReader(nil)
				},
			},

			tlsConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				MinVersion:   tls.VersionTLS12,
			},
		},
		nil
}

func (s *Server) Run(address string) error {
	listener, errListen := net.Listen("tcp", address)
	if errListen != nil {
		return fmt.Errorf("listener start: %w", errListen)
	}
	defer listener.Close()

	listenerTLS := tls.NewListener(listener, s.tlsConfig)
	defer listenerTLS.Close()

	log.Printf(
		"Server %s listening on %s (HTTPS)...\n",

		_Version,
		address,
	)

	for {
		conn, err := listenerTLS.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)

			continue
		}

		go s.handleConnection(
			newConnection(conn),
		)
	}
}

func (s *Server) handleConnection(conn *connection) {
	defer conn.Close()

	// Set a read deadline for idle connections
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	for {
		if errOnTraffic := s.onTraffic(conn); errOnTraffic != nil {
			break
		}

		// Reset the timeout after successful activity
		_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	}
}

func (s *Server) SendStatus(statusCode int) []byte {
	buffer := s.bufferPool.Get().(*bytes.Buffer)
	defer s.bufferPool.Put(buffer)
	buffer.Reset()

	buffer.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode)))
	buffer.WriteString("Content-Length: 0\r\n")
	buffer.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123))) // Add Date header
	buffer.WriteString("\r\n")

	return buffer.Bytes()
}

func (s *Server) SendBody(statusCode int, body string) []byte {
	buffer := s.bufferPool.Get().(*bytes.Buffer)
	defer s.bufferPool.Put(buffer)
	buffer.Reset()

	buffer.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode)))
	buffer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
	buffer.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123))) // Add Date header
	buffer.WriteString("\r\n")
	buffer.WriteString(body)

	return buffer.Bytes()
}

func (s *Server) onTraffic(conn *connection) error {
	data, errRead := conn.Read()
	if errRead != nil {
		if errRead == io.EOF {
			return nil // client closed the connection gracefully
		}

		return errRead // actual read error
	}

	bufReader := s.readerPool.Get().(*bufio.Reader)
	defer func() {
		bufReader.Reset(nil)
		s.readerPool.Put(bufReader)
	}()
	bufReader.Reset(bytes.NewReader(data))

	request, errParse := http.ReadRequest(bufReader)
	if errParse != nil {
		go log.Printf("failed to parse HTTP request: %v", errParse)

		_, _ = conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 0\r\n\r\n"))

		return nil // Don't close the connection on a bad request
	}

	go func() {
		fmt.Printf(
			"IP: %s, Method: %s, Path: %s\n",
			getClientIP(request),
			request.Method,
			request.URL.Path,
		)
	}()

	_, _ = conn.Write(s.SendStatus(http.StatusOK))

	if strings.ToLower(request.Header.Get("Connection")) == "close" {
		return io.EOF // Signal to close the connection
	}

	return nil // Keep the connection open
}
