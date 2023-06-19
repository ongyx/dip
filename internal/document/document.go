package document

import (
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

// Document represents a chunk of Markdown text in UTF-8 as HTML.
type Document struct {
	// Title is the name of the document.
	Title string

	// Timestamp is when the document was last modified with Convert().
	Timestamp time.Time

	sb strings.Builder
}

// String returns the document's content.
func (d *Document) String() string {
	return d.sb.String()
}

// Convert converts the given Markdown text to HTML.
func (d *Document) Convert(text []byte, md goldmark.Markdown) error {
	d.sb.Reset()

	if err := md.Convert(text, &d.sb); err != nil {
		return err
	} else {
		d.Timestamp = time.Now()
	}

	return nil
}
