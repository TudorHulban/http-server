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

	"github.com/TudorHulban/https-server/router"
)

type Server struct {
	bufferPool sync.Pool
	readerPool sync.Pool

	tlsConfig *tls.Config

	router *router.Router
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

			router: router.NewRouter(),
		},
		nil
}

func (s *Server) worker(chConnection chan *connection) {
	for conn := range chConnection {
		defer conn.Close()

		for {
			if errOnTraffic := s.onTraffic(conn); errOnTraffic != nil {
				break
			}

			// Set a read deadline for idle connections
			_ = conn.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))
		}
	}
}

func (s *Server) Run(address string, workerCount, channelCapacity int) error {
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

	// Create a channel for queuing connections
	chConnections := make(chan *connection, channelCapacity)
	defer close(chConnections)

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		go s.worker(chConnections)
	}

	for {
		conn, err := listenerTLS.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)

			continue
		}

		select {
		case chConnections <- newConnection(conn):

		default:
			log.Printf(
				"connection queue is full, rejecting connection from: %s",
				conn.RemoteAddr().String(),
			)
			conn.Close()
		}
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

	if request.Method != "GET" && request.Method != "HEAD" {
		// Read the body only if necessary
		body, err := io.ReadAll(bufReader)
		if err != nil {
			go log.Printf("failed to read request body: %v", err)

			_, _ = conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 0\r\n\r\n"))

			return nil
		}

		// Now request object has complete data including body
		// Process the body here (if needed)
		request.Body = io.NopCloser(bytes.NewReader(body))
	}

	// go func() {
	// 	fmt.Printf(
	// 		"IP: %s, Method: %s, Path: %s\n",
	// 		getClientIP(request),
	// 		request.Method,
	// 		request.URL.Path,
	// 	)
	// }()

	if handler, exists := s.router.FindHandler(request.URL.Path); exists {
		buf := bytes.NewBuffer(nil)
		rw := &responseWriter{w: bufio.NewWriter(buf)}

		handler(rw, request)
		rw.w.Flush()

		_, err := conn.Write(buf.Bytes())
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}

	} else {
		_, _ = conn.Write(s.SendStatus(http.StatusNotFound))
	}

	if strings.ToLower(request.Header.Get("Connection")) == "close" {
		return io.EOF // Signal to close the connection
	}

	return nil // Keep the connection open
}
