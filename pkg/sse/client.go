package sse

import "net/http"

// Client is a HTTP/2 client that supports SSE.
type Client interface {
	http.ResponseWriter
	http.Flusher
}
