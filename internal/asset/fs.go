package asset

import (
	"embed"
	"io/fs"
)

var (
	//go:embed dist
	filesystem embed.FS

	// FS is the default filesystem for the server's assets.
	FS, _ = fs.Sub(filesystem, "dist")
)
