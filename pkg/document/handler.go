package document

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/ongyx/dip/pkg/source"
	"github.com/ongyx/dip/pkg/sse"
)

const (
	// Path for serving application assets.
	assetURL = "__assets"

	// Path for serving SSE events.
	eventURL = "__events"
)

// Handler is a wrapper around a library for serving documents over HTTP.
type Handler struct {
	lib *Library
	log *log.Logger

	sse *sse.Server
}

// NewHandler creates a new handler.
func NewHandler(lib *Library, log *log.Logger) *Handler {
	h := &Handler{
		lib: lib,
		log: log,
		sse: sse.NewServer(""),
	}
	go h.reload()

	return h
}

// Close cleans up the handler.
func (h *Handler) Close() {
	h.sse.Close()
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Clean path and remove leading slash
	p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")

	// If the path is empty, use the root document.
	if p == "" {
		p = source.Root
	}

	d, err := h.lib.Open(p)
	if err != nil {
		// Create the document.
		d, err = h.lib.Create(p)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// Create an SSE stream for the document.
		h.sse.Add(p)
	}

	err = d.Borrow(func(buf []byte) error {
		eu := &url.URL{
			Path: eventURL,
			RawQuery: url.Values{
				"stream": {p},
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
		h.log.Printf("error: serving document %s to %s failed: %s", p, r.Host, err)
	}
}

func (h *Handler) reload() {
	if w, ok := h.lib.Watcher(); ok {
		files, errors := w.Watch()

		for {
			select {
			case f, ok := <-files:
				if !ok {
					return
				}

				h.log.Println("handler: reloading document", f)

				d, err := h.lib.Create(f)
				if err != nil {
					h.log.Printf("error: handler: failed to reload document %s: %s\n", f, err)
					continue
				}

				// Send a reload event over SSE if the stream exists.
				if st := h.sse.Get(f); st != nil {
					d.Borrow(func(buf []byte) error {
						st.Send(&sse.Event{Type: "reload", Data: buf})

						return nil
					})
				}
			case err, ok := <-errors:
				if !ok {
					return
				}

				h.log.Println("error: handler: watcher ran into an error:", err)
			}
		}
	}
}
