package document

import (
	_ "embed"
	"html/template"
	"io"
)

var (
	//go:embed template.html
	rawTemplate string

	documentTemplate = template.Must(template.New("document").Parse(rawTemplate))
)

// Template is the HTML template for serving a document.
// For the template itself, see template.html.
type Template struct {
	// The HTML title to use.
	Title string

	// The directory where assets are served.
	// index.js and index.css must be present in the directory.
	AssetURL string

	// The path where SSE events are served.
	EventURL string

	// The initial content of the HTML page.
	Content template.HTML
}

func (t *Template) Execute(w io.Writer) error {
	return documentTemplate.Execute(w, t)
}
