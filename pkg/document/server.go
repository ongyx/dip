package document

import (
	"log"
	"net/http"
	"net/url"

	"github.com/yuin/goldmark"

	"github.com/ongyx/dip/pkg/asset"
	"github.com/ongyx/dip/pkg/source"
)

// Server serves documents and their accompanying CSS/JS assets over HTTP.
type Server struct {
	handler *Handler
	mux     *http.ServeMux
}

// NewServer creates a new server with a source URL.
func NewServer(u *url.URL, md goldmark.Markdown, log *log.Logger) (*Server, error) {
	src, err := source.New(u)
	if err != nil {
		return nil, err
	}

	h := NewHandler(NewLibrary(src, md), log)

	mux := http.NewServeMux()

	// Serve Markdown documents at the root by default.
	mux.Handle("/", h)

	// Serve assets at the asset path.
	ap := "/" + assetURL + "/"
	mux.Handle(ap, http.StripPrefix(ap, asset.FileServer))

	// Serve SSE events.
	mux.Handle("/"+eventURL, h.sse)

	return &Server{handler: h, mux: mux}, nil
}

// Close closes the server.
func (s *Server) Close() {
	s.handler.Close()
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
