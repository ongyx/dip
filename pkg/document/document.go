package document

import (
	"bytes"
	"io"
	"sync"

	"github.com/yuin/goldmark"
)

// Document is a HTML representation of Markdown text.
type Document struct {
	name string
	md   goldmark.Markdown

	mu  sync.RWMutex
	buf bytes.Buffer
}

// NewDocument creates a new document with the given name and Markdown converter.
func NewDocument(name string, md goldmark.Markdown) *Document {
	return &Document{name: name, md: md}
}

// Name returns the document's name.
func (d *Document) Name() string {
	return d.name
}

// Convert converts the Markdown text in buf into the document as HTML.
func (d *Document) Convert(buf []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.buf.Reset()

	return d.md.Convert(buf, &d.buf)
}

// Borrow temporarily borrows the document's HTML for the lifetime of fn.
// buf must not be mutated or used outside of fn.
func (d *Document) Borrow(fn func(buf []byte) error) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return fn(d.buf.Bytes())
}

// ReaderFrom reads Markdown text from the reader and converts it into the document as HTML.
func (d *Document) ReadFrom(r io.Reader) (n int64, err error) {
	// TODO: perhaps persist this buffer for reading?
	var buf bytes.Buffer

	n, err = buf.ReadFrom(r)
	if err != io.EOF {
		return n, err
	}

	return n, d.Convert(buf.Bytes())
}

// WriteTo writes the document's HTML to the writer.
func (d *Document) WriteTo(w io.Writer) (n int64, err error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.buf.WriteTo(w)
}
