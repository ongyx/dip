package template

import (
	"embed"
	"html/template"
	"io"
)

var (
	//go:embed *.html
	tmplFS embed.FS

	tmpl = template.Must(template.New("template").ParseFS(tmplFS, "*.html"))
)

type Document struct {
	Title     string
	AssetPath string
	Content   template.HTML
}

func (dt *Document) Execute(w io.Writer) error {
	return tmpl.Lookup("document.html").Execute(w, dt)
}
