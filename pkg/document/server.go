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
	mux *http.ServeMux
}

func NewServer(u *url.URL, md goldmark.Markdown, lg *log.Logger) (*Server, error) {
	src, err := source.New(u)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		Lib: NewLibrary(src, md),
		Log: lg,
	}

	mux := http.NewServeMux()

	ap := "/" + assetURL + "/"

	// Serve assets at the asset path.
	mux.Handle(ap, http.StripPrefix(ap, asset.FileServer))

	// Serve Markdown documents at the root by default.
	mux.Handle("/", h)

	return &Server{mux: mux}, nil
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
