package main

import (
	"fmt"
	"strings"
)

type HTTPRequest struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}

func NewHTTPSRequest(raw string) (*HTTPRequest, error) {
	lines := strings.Split(raw, "\r\n")

	if len(lines) < 1 {
		return nil, fmt.Errorf("invalid HTTPS request")
	}

	// Parse the request line
	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) != 3 {
		return nil,
			fmt.Errorf("invalid request line")
	}

	method, path, version := requestLine[0], requestLine[1], requestLine[2]

	headers := make(map[string]string)
	i := 1

	for ; i < len(lines) && lines[i] != ""; i++ {
		header := strings.SplitN(lines[i], ": ", 2)

		if len(header) == 2 {
			headers[header[0]] = header[1]
		}
	}

	body := strings.Join(lines[i+1:], "\r\n")

	return &HTTPRequest{
			Method:  method,
			Path:    path,
			Version: version,
			Headers: headers,
			Body:    body,
		},
		nil
}
