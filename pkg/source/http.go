package source

import (
	"io"
	"net/http"
	"net/url"
	"testing/fstest"
)

// HTTP is a source that serves a URL.
type HTTP struct {
	fstest.MapFS
}

// NewHTTP creates a new source, downloading over HTTP(s) from the given URL.
func NewHTTP(u *url.URL) (Source, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HTTP{
		MapFS: fstest.MapFS{
			".": {Data: data},
		},
	}, nil
}
