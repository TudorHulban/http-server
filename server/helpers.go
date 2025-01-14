package server

import (
	"net/http"
	"strings"
	"time"
)

func getDayOfYear() int {
	now := time.Now()

	return now.YearDay()
}

func getClientIP(request *http.Request) string {
	if forwarded := request.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")

		return strings.TrimSpace(parts[0])
	}

	if realIP := request.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	remoteIP := request.RemoteAddr

	// In some cases RemoteAddr contains port number, so we need to split it
	parts := strings.Split(remoteIP, ":")

	return strings.TrimSpace(parts[0])
}
