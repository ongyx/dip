package web

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/ongyx/dip/pkg/document"
	"github.com/ongyx/dip/pkg/source"
	"github.com/ongyx/dip/pkg/static"
)

const (
	staticPath = "__"
)

// Server serves a library of documents over HTTP.
type Server struct {
	library *document.Library
	logger  *log.Logger
}

// NewServer creates a new server
func NewServer(library *document.Library, logger *log.Logger) *Server {
	library.Watch(logger)

	return &Server{
		library: library,
		logger:  logger,
	}
}

// Mux creates a new multiplexer, given the filesystem to serve static assets from.
// If staticFS is nil, it defaults to static.FS.
func (s *Server) Mux(staticFS fs.FS) *http.ServeMux {
	if staticFS == nil {
		staticFS = static.FS
	}

	mux := http.NewServeMux()

	// serve documents from the root
	mux.Handle("/", s)

	// static resources
	path := "/" + staticPath + "/"
	srv := http.FileServer(http.FS(staticFS))

	mux.Handle(path, http.StripPrefix(path, srv))

	return mux
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	// resolve the URL root to the actual root
	if path == "" {
		path = source.Root
	}

	// bail if the path doesn't have a markdown extension and tbe path isn't root
	if !source.IsMarkdownFile(path) && path != source.Root {
		http.NotFound(w, r)
		return
	}

	// load document for the first time if it does not exist
	if !s.library.Has(path) {
		if err := s.library.Reload(path); err != nil {
			s.logger.Printf("error: loading document %s: %s\n", path, err)

			if pe, ok := err.(*fs.PathError); ok {
				if pe.Err == fs.ErrNotExist {
					http.NotFound(w, r)
				}
			}

			return
		}
	}

	// SAFETY: d is guaranteed to be non-nil
	d := s.library.Get(path)
	d.Borrow(func(content string) {
		dt := &static.DocumentTemplate{
			Title:   d.Title(),
			Static:  staticPath,
			Content: template.HTML(content),
		}

		if err := dt.Execute(w); err != nil {
			s.logger.Printf("error: writing document %s: %s\n", path, err)
		}
	})
}
