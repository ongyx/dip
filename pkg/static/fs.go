package static

import (
	"embed"
)

var (
	// FS is the default filesystem for the server's static resources.
	//
	//go:embed dist
	FS embed.FS
)
