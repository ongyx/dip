package document

import (
	"bytes"
	"io"
	"sync"

	"github.com/yuin/goldmark"
)

// Document is a HTML representation of Markdown text.
type Document struct {
	md goldmark.Markdown

	mu  sync.RWMutex
	buf bytes.Buffer
}

// NewDocument creates a new document with the given Markdown converter.
func NewDocument(md goldmark.Markdown) *Document {
	return &Document{md: md}
}

// ReaderFrom reads Markdown text from the reader and converts it into the document as HTML.
func (d *Document) ReadFrom(r io.Reader) (n int64, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// TODO: perhaps persist this buffer for reading?
	var buf bytes.Buffer

	n, err = buf.ReadFrom(r)
	if err != io.EOF {
		return n, err
	}

	if err = d.md.Convert(buf.Bytes(), &d.buf); err != nil {
		return n, err
	}

	return n, nil
}

// WriteTo writes the document's HTML to the writer.
func (d *Document) WriteTo(w io.Writer) (n int64, err error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.buf.WriteTo(w)
}
