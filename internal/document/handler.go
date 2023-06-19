package document

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ongyx/dip/internal/source"
	"github.com/ongyx/dip/internal/sse"
	tmpl "github.com/ongyx/dip/internal/template"
)

const (
	// AssetPath is the relative path where CSS/JS assets should be served.
	AssetPath = "__"
)

// Handler wraps a library to serve documents over HTTP.
type Handler struct {
	lib *Library
	sse *sse.Server
	log *log.Logger
}

// NewHandler creates a new handler that serves documents from the given library.
func NewHandler(lib *Library, sse *sse.Server, lg *log.Logger) *Handler {
	return &Handler{
		lib: lib,
		sse: sse,
		log: lg,
	}
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, ok := h.clean(r.URL)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// load document for the first time if it does not exist
	if !h.lib.Add(path) {
		if err := h.lib.Reload(path); err != nil {
			h.log.Printf("error: couldn't load document %s: %s\n", path, err)

			if pe, ok := err.(*fs.PathError); ok {
				if pe.Err == fs.ErrNotExist {
					http.NotFound(w, r)
				}
			}

			return
		}

		h.log.Printf("loaded document %s\n", path)
	}

	md, err := h.metadata(path)
	if err != nil {
		h.log.Printf("error: failed to marshal metadata: %s\n", err)
		return
	}

	err = h.lib.Borrow(path, func(d *Document) error {
		// template out document and write it
		dt := &tmpl.Document{
			Title:     d.Title,
			AssetPath: AssetPath,
			Metadata:  template.JS(md),
			Content:   template.HTML(d.String()),
		}

		return dt.Execute(w)
	})

	if err != nil {
		h.log.Printf("error: couldn't write document %s: %s\n", path, err)
		return
	}
}

func (h *Handler) metadata(path string) ([]byte, error) {
	return json.Marshal(map[string]any{
		"path": map[string]string{
			"asset":    AssetPath,
			"document": path,
		},
	})
}

func (h *Handler) clean(uri *url.URL) (string, bool) {
	path := strings.Trim(uri.Path, "/")

	fmt.Println(path)

	// bail if the path doesn't have a markdown extension and tbe path isn't root
	if !(source.IsMarkdownFile(path) || path == "") {
		return "", false
	}

	// resolve the URL root to the actual root
	if path == "" {
		path = h.lib.Root()
	}

	return path, true
}
