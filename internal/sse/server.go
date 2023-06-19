package sse

import (
	"context"
	"errors"
	"net/http"
	"sync"
)

var (
	// ErrSSENotSupported indicates the client does not support server-sent events.
	ErrSSENotSupported = errors.New("client does not support SSE")
)

// Server represents a group of clients that events can be broadcasted to.
// Any errors encountered are sent over the Errors channel.
type Server struct {
	Errors chan error

	events Stream

	mu      sync.RWMutex
	streams map[context.Context]Stream
}

// Server creates a new server.
func NewServer() *Server {
	s := &Server{
		Errors:  make(chan error),
		events:  make(chan *Event),
		streams: make(map[context.Context]Stream),
	}
	go s.send()

	return s
}

// Send broadcasts a message to all clients in the server.
func (s *Server) Send(e *Event) {
	s.events <- e
}

// Recieve adds a client to the server and waits for sent events.
// If the client does not support SSE, ErrSSENotSupported is returned.
func (s *Server) Recieve(w http.ResponseWriter, r *http.Request) error {
	c, ok := w.(Client)
	if !ok {
		return ErrSSENotSupported
	}

	h := c.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")

	stream := make(Stream)

	s.streams[r.Context()] = stream

	for e := range stream {
		if _, err := e.WriteTo(c); err != nil {
			s.Errors <- err
		} else {
			c.Flush()
		}
	}

	return nil
}

// Close closes the server, dropping all events.
func (s *Server) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	close(s.events)
	for _, stream := range s.streams {
		close(stream)
	}

	close(s.Errors)
}

func (s *Server) send() {
	for e := range s.events {
		s.broadcast(e)
	}
}

func (s *Server) broadcast(e *Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for ctx, stream := range s.streams {
		select {
		case <-ctx.Done():
			// client disconnected, so remove them.
			delete(s.streams, ctx)

		default:
			stream <- e
		}
	}
}
