package server

import (
	"bufio"
	"net/http"
)

type responseWriter struct {
	w          *bufio.Writer
	statusCode int
	header     http.Header
}

// Header returns the header map.
func (rw *responseWriter) Header() http.Header {
	return rw.header
}

// Write writes the data to the underlying writer.
func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.w.Write(b)
}

// WriteHeader sets the status code for the response.
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}
