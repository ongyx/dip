package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ongyx/dip/internal/asset"
	"github.com/ongyx/dip/internal/document"
	"github.com/ongyx/dip/internal/sse"
)

// Server serves documents and assets.
type Server struct {
	lib *document.Library
	log *log.Logger

	sse *sse.Server
	mux *http.ServeMux
}

// NewServer creates a new server, serving documents from the given library.
func NewServer(lib *document.Library, lg *log.Logger) *Server {
	s := &Server{
		lib: lib,
		log: lg,

		sse: sse.NewServer(),
		mux: http.NewServeMux(),
	}
	go s.watch(lib)

	// documents
	s.mux.Handle("/", document.NewHandler(lib, s.sse, lg))

	// embedded assets
	ap := "/" + document.AssetPath + "/"
	s.mux.Handle(ap, http.StripPrefix(ap, asset.FileServer))

	// server-sent events
	s.mux.HandleFunc(ap+"events", s.handleSSE)

	return s
}

// Close closes the server.
func (s *Server) Close() error {
	s.sse.Close()
	return s.lib.Close()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	if err := s.sse.Recieve(w, r); err != nil {
		s.log.Printf("error: sse: %s\n", err)
	}
}

func (s *Server) watch(lib *document.Library) {
	files, errors := lib.Watch()

	for {
		select {
		case file, ok := <-files:
			if !ok {
				return
			}

			s.log.Printf("reloaded document %s\n", file)

			// publish the content of reloaded files to the SSE server
			go s.send(file)
		case err, ok := <-errors:
			if !ok {
				return
			}

			s.log.Printf("error: watcher: %s\n", err)
		}
	}
}

func (s *Server) send(file string) {
	s.lib.Borrow(file, func(d *document.Document) error {
		msg := Message{
			Content:   d.String(),
			Timestamp: d.Timestamp.Unix(),
		}

		b, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		var e sse.Event

		s.sse.Send(e.Type(file).Data(b))

		return nil
	})
}
