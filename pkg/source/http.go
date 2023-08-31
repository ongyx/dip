package source

import (
	"net/http"
	"net/url"
)

// HTTP is a source that serves a URL.
type HTTP struct {
	*VirtualFS
}

// NewHTTP creates a new source, downloading over HTTP(s) from the given URL.
func NewHTTP(u *url.URL) (Source, error) {
	url := u.String()
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	vf := NewVirtualFile(url)

	if _, err := vf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	return &HTTP{NewVirtualFS(vf)}, nil
}
