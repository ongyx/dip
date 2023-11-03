package sse

import (
	"bytes"
	"context"
	"net/http"
	"sync"
)

// Stream represents an event stream with several clients.
// Any errors encountered are sent over the Errors channel.
type Stream struct {
	events chan *Event

	mu      sync.RWMutex
	clients map[context.Context]chan []byte
}

// NewStream creates a new stream.
func NewStream() *Stream {
	s := &Stream{
		events:  make(chan *Event),
		clients: make(map[context.Context]chan []byte),
	}

	go func() {
		// Cache for marshalling events.
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
// If the client does not support SSE, a 500 status code will be sent.
func (s *Stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// Prepare headers for SSE streaming.
	// See https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events.
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")

	// Create a channel to receive marshaled events from and add the client to the pool.
	marshaled := make(chan []byte)
	s.clients[r.Context()] = marshaled

	// Signal to the client that connection setup is complete.
	w.WriteHeader(http.StatusOK)
	f.Flush()

	for m := range marshaled {
		w.Write(m)
		f.Flush()
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
}

func (s *Stream) send(e *Event, buf *bytes.Buffer) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	buf.Reset()
	e.Marshal(buf)

	for ctx, marshaled := range s.clients {
		select {
		case <-ctx.Done():
			// Client disconnected, so remove them.
			delete(s.clients, ctx)

		default:
			marshaled <- buf.Bytes()
		}
	}
}
