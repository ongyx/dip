package source

import (
	"io/fs"
	"net/url"
	"os"
)

const (
	// Root represents the root document path in a source.
	// Certain sources only serve a single Markdown file and thus may check for this value to serve that file.
	Root = "."
)

// Source is a filesystem serving Markdown files.
type Source interface {
	fs.FS
}

// New creates a source from the URL.
func New(u *url.URL) (Source, error) {
	return Get(u.Scheme)(u)
}

// Parse parses the path into a source.
//
// The path is parsed in order as follows:
//
// * If the path is already a URI, it is passed verbatim to New().
//
// * If the path is a dash ('-'), standard input is used.
//
// * Otherwise, stat the path to check if the path is a file or directory.
func Parse(path string) (*url.URL, error) {
	if u, err := url.ParseRequestURI(path); err == nil {
		return u, nil
	}

	u := &url.URL{Path: path}

	// Read from standard input if given a dash.
	if path == "-" {
		u.Scheme = "stdin"
	} else {
		stat, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		// Check if the path is a file or directory.
		if stat.IsDir() {
			u.Scheme = "dir"
		} else {
			u.Scheme = "file"
		}
	}

	return u, nil
}

// Must unwraps a (value, error) return to just the value.
// If err is not nil, this panics.
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}
