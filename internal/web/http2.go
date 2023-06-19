package web

import (
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// H2Cify wraps the handler in a HTTP/2 Cleartext server to allow HTTP/2 connections (and SSE) without HTTPS.
func H2Cify(h http.Handler) http.Handler {
	return h2c.NewHandler(h, &http2.Server{})
}
