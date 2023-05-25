package static

import (
	"embed"
	"html/template"
	"io"
)

var (
	//go:embed template
	tmplFS embed.FS

	tmpl = template.Must(template.New("template").ParseFS(tmplFS, "template/*.html"))
)

type DocumentTemplate struct {
	Title   string
	Static  string
	Content template.HTML
}

func (dt *DocumentTemplate) Execute(w io.Writer) error {
	return tmpl.Lookup("document.html").Execute(w, dt)
}
