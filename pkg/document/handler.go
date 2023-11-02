package document

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/ongyx/dip/pkg/source"
)

const (
	// Path for serving application assets.
	assetURL = "__assets"

	// Path for serving SSE events.
	eventURL = "__events"
)

// Handler is a wrapper around a library for serving documents over HTTP.
type Handler struct {
	Lib *Library
	Log *log.Logger
}

// ServeHTTP implements http.Handler.
func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := path.Clean(r.URL.Path)

	if p == "/" {
		p = source.Root
	}

	d, err := s.Lib.Open(p)
	if err != nil {
		// Create the document.
		d, err = s.Lib.Create(p)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}

	err = d.Borrow(func(buf []byte) error {
		eu := &url.URL{
			Path: eventURL,
			RawQuery: url.Values{
				"path": {p},
			}.Encode(),
		}

		t := Template{
			Title:    "",
			AssetURL: assetURL,
			EventURL: eu.String(),
			Content:  template.HTML(buf),
		}
		if err := t.Execute(w); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.Log.Printf("error: serving document %s to %s failed: %s", p, r.Host, err)
	}
}
