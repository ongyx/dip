package sse

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"sync"
)

var (
	// ErrSSENotSupported indicates the client does not support stream-sent events.
	ErrSSENotSupported = errors.New("client does not support SSE")
)

// Stream represents an event stream with several clients.
// Any errors encountered are sent over the Errors channel.
type Stream struct {
	Errors chan error

	events chan *Event

	mu      sync.RWMutex
	clients map[context.Context]chan []byte
}

// Stream creates a new stream.
func NewStream() *Stream {
	s := &Stream{
		Errors:  make(chan error),
		events:  make(chan *Event),
		clients: make(map[context.Context]chan []byte),
	}

	go func() {
		// Cache for marshalling events
		var buf bytes.Buffer

		for e := range s.events {
			s.send(e, &buf)
		}
	}()

	return s
}

// Send sends an event to all clients connected to the stream.
func (s *Stream) Send(e *Event) {
	s.events <- e
}

// ServeHTTP connects a client to the stream and waits for sent events.
// If the client does not support SSE, sse.Error is sent over the error channel with ErrSSENotSupported.
func (s *Stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, ok := w.(Client)
	if !ok {
		s.Errors <- &Error{error: ErrSSENotSupported, Request: r}
	}

	h := c.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")

	marshaled := make(chan []byte)

	s.clients[r.Context()] = marshaled

	for m := range marshaled {
		if _, err := c.Write(m); err != nil {
			s.Errors <- err
		} else {
			c.Flush()
		}
	}
}

// Close closes the stream, dropping all events.
func (s *Stream) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	close(s.events)
	for _, ev := range s.clients {
		close(ev)
	}

	close(s.Errors)
}

func (s *Stream) send(e *Event, buf *bytes.Buffer) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	buf.Reset()
	e.Marshal(buf)

	for ctx, marshaled := range s.clients {
		select {
		case <-ctx.Done():
			// client disconnected, so remove them.
			delete(s.clients, ctx)

		default:
			marshaled <- buf.Bytes()
		}
	}
}
