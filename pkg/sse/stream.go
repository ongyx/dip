package sse

import (
	"bytes"
	"sync"
)

// Stream represents a group of clients that messages can be broadcasted to.
type Stream struct {
	Errors chan error

	messages chan Message

	mu      sync.RWMutex
	clients map[Client]bool
}

// Stream creates a new stream.
func NewStream() *Stream {
	s := &Stream{
		Errors:   make(chan error),
		messages: make(chan Message),
		clients:  make(map[Client]bool),
	}
	go s.send()

	return s
}

// Send broadcasts a message to all clients in the stream.
func (s *Stream) Send(m Message) {
	s.messages <- m
}

// Add adds a client to the stream.
func (s *Stream) Add(c Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[c] = true
}

// Close closes the stream, dropping all messages.
func (s *Stream) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	close(s.messages)
	close(s.Errors)
}

func (s *Stream) send() {
	var buf bytes.Buffer

	for msg := range s.messages {
		msg.Marshal(&buf)
		s.broadcast(&buf)
	}
}

func (s *Stream) broadcast(buf *bytes.Buffer) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for client := range s.clients {
		select {
		case <-client.context.Done():
			// client disconnected, so remove them.
			delete(s.clients, client)

		default:
			// send buffer directly to client
			// this avoids having to marshal the message for each client
			if _, err := client.send(buf); err != nil {
				s.Errors <- err
			}
		}
	}
}
