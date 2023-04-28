package dip

import (
	"bytes"
	"sync"
)

// Index is an index of documents mapped by paths.
type Index struct {
	mu        sync.RWMutex
	documents map[string]*bytes.Buffer
}

// Index creates a new index.
func NewIndex() *Index {
	return &Index{documents: make(map[string]*bytes.Buffer)}
}

// Cache caches the templated document in the index.
func (idx *Index) Cache(path string, doc *Document) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	b, ok := idx.documents[path]
	if !ok {
		b = new(bytes.Buffer)
		idx.documents[path] = b
	}

	b.Reset()

	return doc.Execute(b)
}

// Borrow borrows the document buffer at path, or nil if the path does not exist.
// The buffer must not be retained outside of fn.
func (idx *Index) Borrow(path string, fn func(buf []byte)) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var buf []byte
	if b, ok := idx.documents[path]; ok {
		buf = b.Bytes()
	}

	fn(buf)
}

// Has checks if a document exists at path.
func (idx *Index) Has(path string) bool {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	_, ok := idx.documents[path]
	return ok
}
