package web

import (
	"io/fs"
	"net/http"

	"github.com/r3labs/sse/v2"
)

type Endpoint struct {
	SSE *sse.Server

	path  string
	asset fs.FS
}

func NewEndpoint(path string, asset fs.FS) *Endpoint {
	return &Endpoint{
		SSE:   sse.New(),
		path:  path,
		asset: asset,
	}
}

func (ep *Endpoint) Mux() *http.ServeMux {
	mux := http.NewServeMux()

	// serve assets
	ap := "/" + ep.path + "/"
	as := http.FileServer(http.FS(ep.asset))
	mux.Handle(ap, http.StripPrefix(ap, as))

	// serve http2 server-sent events
	mux.Handle("/events", ep.SSE)

	return mux
}
