package web

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ongyx/dip/internal/document"
	"github.com/ongyx/dip/internal/source"
	tmpl "github.com/ongyx/dip/internal/template"
)

type documentHandler struct {
	library *document.Library
	logger  *log.Logger
}

func (h *documentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.clean(r.URL)

	// bail if the path doesn't have a markdown extension and tbe path isn't root
	if !source.IsMarkdownFile(path) && path != source.Root {
		http.NotFound(w, r)
		return
	}

	// load document for the first time if it does not exist
	if err := h.load(path); err != nil {
		h.logger.Printf("error: loading document %s: %s\n", path, err)

		if pe, ok := err.(*fs.PathError); ok {
			if pe.Err == fs.ErrNotExist {
				http.NotFound(w, r)
			}
			return
		}

	}

	// SAFETY: d is guaranteed to be non-nil
	d := h.library.Get(path)
	d.Borrow(func(content string) {
		dt := &tmpl.Document{
			Title:     d.Title(),
			AssetPath: AssetPath,
			Content:   template.HTML(content),
		}

		if err := dt.Execute(w); err != nil {
			h.logger.Printf("error: writing document %s: %s\n", path, err)
		}
	})
}

func (h *documentHandler) load(path string) error {
	if !h.library.Has(path) {
		return h.library.Reload(path)
	}

	return nil
}

func (h *documentHandler) clean(uri *url.URL) string {
	path := strings.Trim(uri.Path, "/")

	// resolve the URL root to the actual root
	if path == "" {
		path = source.Root
	}

	return path
}
