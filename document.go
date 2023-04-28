package dip

import (
	"html/template"
	"io"
	"strings"

	"github.com/yuin/goldmark"
)

// Document represents a Markdown document converted into HTML.
type Document struct {
	Title   string
	Static  string
	Content template.HTML

	content strings.Builder
}

// Convert converts Markdown text to HTML.
func (d *Document) Convert(text []byte, md goldmark.Markdown) error {
	d.content.Reset()

	if err := md.Convert(text, &d.content); err != nil {
		return err
	}

	d.Content = template.HTML(d.content.String())

	return nil
}

// Execute executes the document template.
func (d *Document) Execute(w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "document.html", d)
}
