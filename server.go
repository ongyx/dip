package dip

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

// Server serves documents over HTTP.
type Server struct {
	library *Library

	log *log.Logger
}

// NewServer creates a new server.
// Errors are logged to the given logger.
func NewServer(library *Library, logger *log.Logger) *Server {
	go func() {
		for event := range library.Watch() {
			title := library.Title(event.Path)

			if event.Error != nil {
				logger.Printf("error: couldn't reload document %s: %s\n", title, event.Error)
			} else {
				logger.Printf("reloaded document %s\n", title)
			}
		}
	}()

	return &Server{library: library, log: logger}
}

// Mux creates a new multiplexer and registers the server as a handler.
// If static is nil, dip.Static is used as the default filesystem.
func (s *Server) Mux(static fs.FS) *http.ServeMux {
	if static == nil {
		static = Static
	}

	mux := http.NewServeMux()

	// serve documents from the root
	mux.Handle("/", s)

	// static resources
	staticPath := fmt.Sprintf("/%s/", strings.Trim(s.library.Static, "/"))
	fileServer := http.FileServer(http.FS(static))
	mux.Handle(staticPath, http.StripPrefix(staticPath, fileServer))

	return mux
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")

	if !s.library.Has(path) {
		if err := s.library.Reload(path); err != nil {
			s.log.Printf("error: couldn't init document at %s: %s\n", s.library.Title(path), err)
		}
	}

	var written int
	var err error

	s.library.Borrow(path, func(buf []byte) {
		written, err = w.Write(buf)
	})

	if err != nil {
		s.log.Printf("error: wrote only %d bytes to response: %s\n", written, err)
	}
}
