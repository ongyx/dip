package asset

import (
	"embed"
	"net/http"
)

var (
	// FS is the default filesystem for the server's assets.
	//go:embed dist
	FS embed.FS

	// FileServer handles asset requests.
	FileServer = http.FileServer(http.FS(FS))
)
