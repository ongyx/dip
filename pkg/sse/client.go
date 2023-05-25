package sse

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
)

var (
	// ErrSSENotSupported indicates the HTTP handler does not support server-sent events.
	ErrSSENotSupported = errors.New("response writer does not support SSE")
)

// Client represents a HTTP/2 client that can receive server-sent events.
type Client struct {
	writer  http.ResponseWriter
	flusher http.Flusher
	context context.Context
}

// NewClient creates a new client.
// If the client does not support SSE, ErrSSENotSupported is returned as the error.
func NewClient(w http.ResponseWriter, r *http.Request) (Client, error) {
	f, ok := w.(http.Flusher)
	if !ok {
		return Client{}, ErrSSENotSupported
	}

	h := w.Header()

	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-store")

	return Client{writer: w, flusher: f, context: r.Context()}, nil
}

// Send sends the message to the client, returning the number of message bytes written and the error encountered, if any.
func (c Client) Send(msg Message) (int64, error) {
	var buf bytes.Buffer
	msg.Marshal(&buf)

	return c.send(&buf)
}

func (c Client) send(buf *bytes.Buffer) (int64, error) {
	defer c.flusher.Flush()

	return io.Copy(c.writer, buf)
}
