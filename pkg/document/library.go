package document

import (
	"io/fs"
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"github.com/ongyx/dip/pkg/source"
)

var (
	// DefaultMarkdown is the default converter for Markdown to HTML.
	DefaultMarkdown = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
)

// Library is a collection of Markdown documents.
type Library struct {
	src source.Source
	md  goldmark.Markdown

	mu   sync.RWMutex
	docs map[string]*Document
}

// NewLibrary creates a new library using the given source and Markdown converter.
// If md is nil, DefaultMarkdown is used.
func NewLibrary(src source.Source, md goldmark.Markdown) *Library {
	if md == nil {
		md = DefaultMarkdown
	}

	return &Library{
		src:  src,
		md:   md,
		docs: make(map[string]*Document),
	}
}

// Open reads the document with name.
// If the document does not exist, ok is false.
func (l *Library) Open(name string) (d *Document, ok bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	d, ok = l.docs[name]
	return
}

// Create reads a document from the library's source.
// If the document exists, the existing document's content is reloaded.
func (l *Library) Create(name string) (*Document, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	d, ok := l.docs[name]
	if !ok {
		d = NewDocument(name, l.md)
		l.docs[name] = d
	}

	b, err := fs.ReadFile(l.src, name)
	if err != nil {
		return nil, err
	}

	return d, d.Convert(b)
}

// Remove deletes a document from the library.
// If the document does not exist, ok is false.
func (l *Library) Remove(name string) (ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, ok = l.docs[name]
	if ok {
		delete(l.docs, name)
	}

	return ok
}

// Watcher returns the library's source as a Watcher.
// If the source does not support file watching, ok is false.
func (l *Library) Watcher() (w source.Watcher, ok bool) {
	w, ok = l.src.(source.Watcher)
	return
}
