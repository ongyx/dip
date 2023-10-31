package source

import (
	"net/url"
)

var (
	handlers = make(map[string]Handler)
)

func init() {
	// Default handlers.
	Register("dir", NewDirectory)
	Register("file", NewFile)
	Register("http", NewHTTP)
	Register("stdin", NewStdin)
}

// Handler is a function that creates a source from a URL.
type Handler func(*url.URL) (Source, error)

// Register adds a handler for a URL scheme.
func Register(scheme string, handler Handler) {
	handlers[scheme] = handler
}

// Get returns a handler for the URL scheme.
func Get(scheme string) Handler {
	return handlers[scheme]
}
