package sse

import (
	"crypto/tls"
	"net"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// H2CWrap wraps a HTTP handler to support HTTP2 cleartext.
// Note that net/http.Server by itself enables HTTP2 automatically if using a TLS certificate.
func H2CWrap(h http.Handler) http.Handler {
	return h2c.NewHandler(h, &http2.Server{})
}

// H2CTransport creates a transport that supports HTTP2 cleartext.
func H2CTransport() http.RoundTripper {
	return &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}
}
