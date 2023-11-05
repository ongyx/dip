package document

import (
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/yuin/goldmark"

	"github.com/ongyx/dip/pkg/asset"
	"github.com/ongyx/dip/pkg/source"
	"github.com/ongyx/dip/pkg/sse"
)

const (
	// Path for serving application assets.
	assetURL = "/__assets"

	// Path for serving SSE events.
	eventURL = "/__events"
)

// Server serves documents and their accompanying CSS/JS assets over HTTP.
type Server struct {
	sourceURL *url.URL

	lib *Library
	sse *sse.Server
	mux *http.ServeMux

	log *log.Logger
}

// NewServer creates a new server with a source URL.
func NewServer(url *url.URL, md goldmark.Markdown) (*Server, error) {
	src, err := source.New(url)
	if err != nil {
		return nil, err
	}

	s := &Server{
		sourceURL: url,

		lib: NewLibrary(src, md),
		sse: sse.NewServer(""),
		mux: http.NewServeMux(),

		log: log.Default(),
	}

	// Serve Markdown documents at the root by default.
	s.mux.HandleFunc("/", s.serveDocument)

	// Serve assets at the asset path.
	// A trailing slash is needed to indicate any subdirectory requests should go to the file server.
	au := assetURL + "/"
	s.mux.Handle(au, http.StripPrefix(au, asset.FileServer))

	// Serve SSE events.
	s.mux.Handle(eventURL, s.sse)

	return s, nil
}

// SetLogger sets the server's logger.
// By default, the global default logger from log.Default() is used.
func (s *Server) SetLogger(log *log.Logger) {
	s.log = log
}

// Close closes the server.
func (s *Server) Close() {
	s.sse.Close()
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) serveDocument(w http.ResponseWriter, r *http.Request) {
	// Clean path and remove leading slash
	p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")

	// If the path is empty, use the root document.
	if p == "" {
		p = source.Root
	}

	t, err := s.prepareDocument(p)
	if err != nil {
		s.log.Printf("error: preparing document %s failed: %s\n", p, err)

		// Respond with 404 if the document does not exist, otherwise 500.
		if errors.Is(err, fs.ErrNotExist) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "document could not be prepared", http.StatusInternalServerError)
		}

		return
	}

	// Write the document template.
	if err := t.Execute(w); err != nil {
		s.log.Printf("error: serving document %s to %s failed: %s\n", p, r.Host, err)
	}
}

func (s *Server) prepareDocument(name string) (*Template, error) {
	d, err := s.getDocument(name)
	if err != nil {
		return nil, err
	}

	var t *Template

	d.Borrow(func(buf []byte) error {
		ti := d.Name()
		if name == source.Root {
			ti = s.sourceURL.String()
		}

		eu := &url.URL{
			Path: eventURL,
			RawQuery: url.Values{
				"stream": {name},
			}.Encode(),
		}

		t = &Template{
			Title:    ti,
			AssetURL: assetURL,
			EventURL: eu.String(),
			Content:  template.HTML(buf),
		}

		return nil
	})

	return t, nil
}

func (s *Server) getDocument(name string) (*Document, error) {
	if d, ok := s.lib.Open(name); ok {
		return d, nil
	}

	// Create the document.
	d, err := s.lib.Create(name)
	if err != nil {
		return nil, err
	}

	// Create an SSE stream for the document.
	s.sse.Add(name)

	return d, nil
}

func (s *Server) watch() {
	if w, ok := s.lib.Watcher(); ok {
		files, errors := w.Watch()

		for {
			select {
			case f, ok := <-files:
				if !ok {
					return
				}

				s.log.Println("reloading document", f)

				s.reload(f)

			case err, ok := <-errors:
				if !ok {
					return
				}

				s.log.Println("error: watcher:", err)
			}
		}
	}
}

func (s *Server) reload(name string) {
	d, err := s.lib.Create(name)
	if err != nil {
		s.log.Printf("error: handler: failed to reload document %s: %s\n", name, err)
		return
	}

	// Send a reload event over SSE if the stream exists.
	if st := s.sse.Get(name); st != nil {
		d.Borrow(func(buf []byte) error {
			st.Send(&sse.Event{Type: "reload", Data: buf})

			return nil
		})
	}

}
