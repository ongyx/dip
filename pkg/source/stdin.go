package source

import (
	"io"
	"net/url"
	"os"
	"testing/fstest"
)

const (
	stdinFilename = "(stdin)"
)

// Stdin is a source that serves standard input.
type Stdin struct {
	fstest.MapFS
}

// NewStdin creates a new standard input source.
func NewStdin(_ *url.URL) (Source, error) {
	// This reads all data until the user presses Ctrl+D.
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	return &Stdin{
		MapFS: fstest.MapFS{
			".": {Data: data},
		},
	}, nil
}
