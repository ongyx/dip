package document

import (
	"strings"
	"sync"

	"github.com/yuin/goldmark"
)

// Document represents a chunk of Markdown text in UTF-8 as HTML.
type Document struct {
	title string

	mu sync.RWMutex
	sb strings.Builder
}

// Title returns the document title, if any.
func (d *Document) Title() string {
	return d.title
}

// Borrow allows access to the underlying document content.
// buf should not be retained outside of fn.
func (d *Document) Borrow(fn func(content string)) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	fn(d.sb.String())
}

// Convert converts the given Markdown text to HTML.
func (d *Document) Convert(text []byte, md goldmark.Markdown) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.sb.Reset()

	if err := md.Convert(text, &d.sb); err != nil {
		return err
	}

	return nil
}
