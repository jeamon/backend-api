package web

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
)

// generateRequestID provides a random uid to trace a given request.
func generateRequestID() string {
	id, _ := uuid.NewV4()
	return id.String()
}

// getIP helps find the source IP of the caller.
func getIP(r *http.Request) string {
	// Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	// Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP = net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}

	// Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	return "" // "No valid ip found"
}

// readBody convert the body into a string.
// Ref https://github.com/gin-gonic/gin/issues/961#issuecomment-312504339
func readBody(reader io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	s := buf.String()
	return s, nil
}
