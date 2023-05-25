package static

//go:generate curl https://raw.githubusercontent.com/sindresorhus/github-markdown-css/main/github-markdown.css -o css/github.css

import (
	"embed"
)

var (
	// FS is the default filesystem for the server's static resources.
	//
	//go:embed css
	FS embed.FS
)
