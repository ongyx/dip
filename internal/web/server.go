package web

import (
	"log"
	"net/http"

	"github.com/ongyx/dip/internal/asset"
	"github.com/ongyx/dip/internal/document"
	"github.com/r3labs/sse/v2"
)

const (
	// AssetPath is the path where assets are served.
	AssetPath = "__"
)

type Server struct {
	mux *http.ServeMux
	sse *sse.Server
}

func NewServer(library *document.Library, logger *log.Logger) *Server {
	s := &Server{
		mux: http.NewServeMux(),
		sse: sse.New(),
	}

	// documents
	s.mux.Handle("/", &documentHandler{library: library, logger: logger})

	ap := "/" + AssetPath + "/"

	// embedded assets
	s.mux.Handle(ap, http.StripPrefix(ap, asset.FileServer))

	// http2 server-sent events
	s.mux.Handle(ap+"events", s.sse)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
