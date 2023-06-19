package source

import (
	"errors"
	"io/fs"
	"net/url"
)

var (
	// ErrSourceNotFound is returned by New() when the source was not found by the scheme.
	ErrSourceNotFound = errors.New("source not found")
)

// Source represents a filesystem with Markdown files.
type Source interface {
	fs.FS

	// Root returns the path that should be opened in lieu of the root path (i.e., / redirects to README.md).
	Root() string

	// Watch watches the filesystem for changes to Markdown files, sending their paths over the channel files.
	// Errors should be sent over the channel errors.
	Watch(files chan<- string, errors chan<- error)

	// Close closes the source.
	Close() error
}

// New creates a new source given a URI in the format '{scheme}://{path}', i.e. 'file://path/to/markdown/file.md'.
// If the scheme is not found, ErrSourceNotFound is returned.
func New(uri string) (Source, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if fn, ok := sources[u.Scheme]; ok {
		return fn(u.Path)
	}

	return nil, ErrSourceNotFound
}
