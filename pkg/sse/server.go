package sse

import "net/http"

// Server is a stream multiplexer, allowing clients to receive events from a specific stream.
type Server struct {
	defaultStream string
	streams       map[string]*Stream
}

// NewServer creates a new server.
//
// If defaultStream is not empty, a default stream is created for clients who do not specify a stream to connect to.
func NewServer(defaultStream string) *Server {
	s := &Server{
		defaultStream: defaultStream,
		streams:       make(map[string]*Stream),
	}

	if defaultStream != "" {
		s.Add(defaultStream)
	}

	return s
}

// Add creates a new stream in the server.
// If the stream already exists, the existing one is returned.
func (s *Server) Add(name string) *Stream {
	if st, ok := s.streams[name]; ok {
		return st
	}

	st := NewStream()
	s.streams[name] = st

	return st
}

// Get returns an existing stream by name in the server.
// If the stream does not exist, nil is returned.
func (s *Server) Get(name string) *Stream {
	return s.streams[name]
}

// Remove deletes the stream by name from the server and closes it.
// If the stream does not exist, this is a no-op.
func (s *Server) Remove(name string) {
	if st, ok := s.streams[name]; ok {
		st.Close()
		delete(s.streams, name)
	}
}

// Send is a convenience function for Get(name).Send(event).
// The stream must have already been added beforehand; if Get(name) is nil, a panic occurs.
func (s *Server) Send(name string, e *Event) {
	st := s.Get(name)
	if st == nil {
		// More descriptive panic instead of just 'nil pointer deference'.
		panic("server: stream does not exist: " + name)
	}

	st.Send(e)
}

// ServeHTTP connects a client to a stream in the server depending on the 'stream' query parameter.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("stream")
	if name == "" {
		if s.defaultStream != "" {
			name = s.defaultStream
		} else {
			http.Error(w, "SSE event stream not specified", http.StatusBadRequest)
			return
		}
	}

	st := s.Get(name)
	if st == nil {
		http.Error(w, "SSE event stream not found", http.StatusNotFound)
		return
	}

	st.ServeHTTP(w, r)
}
