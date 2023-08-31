package document

import (
	"log"
	"net/http"
	"path"

	"github.com/ongyx/dip/pkg/source"
)

// Server is a HTTP handler that serves Markdown documents.
type Server struct {
	lib *Library
	log *log.Logger
}

// NewServer creates a new document server.
func NewServer(lib *Library, log *log.Logger) *Server {
	return &Server{
		lib: lib,
		log: log,
	}
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := path.Clean(r.URL.Path)
	if p == "/" {
		p = source.Root
	}

	d, err := s.lib.Open(p)
	if err != nil {
		// Create the document.
		d, err = s.lib.Create(p)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}

	if _, err := d.WriteTo(w); err != nil {
		s.log.Printf("server: write failed for %s - %s\n", r.Host, err)
		return
	}
}
