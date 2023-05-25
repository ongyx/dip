package static

//go:generate npx esbuild ../../src/index.js --bundle --outdir=dist

import (
	"embed"
)

var (
	// FS is the default filesystem for the server's static resources.
	//
	//go:embed dist
	FS embed.FS
)
