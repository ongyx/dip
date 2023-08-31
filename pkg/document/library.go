package document

import (
	"errors"
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

	// ErrDocumentNotFound means that a document could not be found in a library.
	ErrDocumentNotFound = errors.New("document not found")
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

// NewLibraryFromPath creates a new library by sourcing a path.
func NewLibraryFromPath(path string, md goldmark.Markdown) (*Library, error) {
	src, err := source.New(path)
	if err != nil {
		return nil, err
	}

	return NewLibrary(src, md), nil
}

// Open reads the document at path.
func (l *Library) Open(path string) (*Document, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if d, ok := l.docs[path]; ok {
		return d, nil
	}

	return nil, ErrDocumentNotFound
}

// Create reads a document from the library's source.
// If the document exists, the existing document's content is reloaded.
func (l *Library) Create(path string) (*Document, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	d, ok := l.docs[path]
	if !ok {
		d = NewDocument(l.md)
	}

	b, err := fs.ReadFile(l.src, path)
	if err != nil {
		return nil, err
	}

	return d, d.Convert(b)
}

// Remove deletes a document from the library.
// If the document does not exist, ErrDocumentNotFound is returned.
func (l *Library) Remove(path string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.docs[path]; !ok {
		return ErrDocumentNotFound
	}

	delete(l.docs, path)
	return nil
}
