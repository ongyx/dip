package sse

import (
	"fmt"
	"net/http"
)

// Error represents an SSE-specific error.
type Error struct {
	error

	Request *http.Request
}

// Error returns an error string with the client's IP address.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.error, e.Request.RemoteAddr)
}
