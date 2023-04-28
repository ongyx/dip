package dip

//go:generate curl https://raw.githubusercontent.com/sindresorhus/github-markdown-css/main/github-markdown.css -o static/css/github.css

import (
	"embed"
	"io/fs"
)

var (
	//go:embed static
	static embed.FS

	// Static is the default filesystem for the server's static resources.
	Static, _ = fs.Sub(static, "static")
)
