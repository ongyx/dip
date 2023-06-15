package document

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/ongyx/dip/pkg/source"
	tmpl "github.com/ongyx/dip/pkg/template"
)

// Handler serves a library of documents over HTTP.
type Handler struct {
	endpoint string
	library  *Library
	logger   *log.Logger
}

func NewHandler(endpoint string, library *Library, logger *log.Logger) *Handler {
	library.Watch(logger)

	return &Handler{
		endpoint: endpoint,
		library:  library,
		logger:   logger,
	}
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	if !h.library.Has(path) {
		if err := h.library.Reload(path); err != nil {
			h.logger.Printf("error: loading document %s: %s\n", path, err)

			if pe, ok := err.(*fs.PathError); ok {
				if pe.Err == fs.ErrNotExist {
					http.NotFound(w, r)
				}
			}

			return
		}
	}

	// SAFETY: d is guaranteed to be non-nil
	d := h.library.Get(path)
	d.Borrow(func(content string) {
		dt := &tmpl.Document{
			Title:    d.Title(),
			Endpoint: h.endpoint,
			Content:  template.HTML(content),
		}

		if err := dt.Execute(w); err != nil {
			h.logger.Printf("error: writing document %s: %s\n", path, err)
		}
	})
}
