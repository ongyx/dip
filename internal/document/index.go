package document

import (
	"sync"
)

// Index is a concurrent-safe map of paths to documents.
type Index struct {
	mu   sync.RWMutex
	docs map[string]*document
}

// NewIndex creates a new index.
func NewIndex() *Index {
	return &Index{docs: make(map[string]*document)}
}

// Add adds a new document to the index.
// exists is true if the path was already added.
func (i *Index) Add(path string) (exists bool) {
	i.mu.Lock()
	defer i.mu.Unlock()

	_, ok := i.docs[path]
	if !ok {
		i.docs[path] = &document{doc: &Document{}}
	}

	return ok
}

// Has checks if the index has a document.
func (i *Index) Has(path string) (exists bool) {
	_, ok := i.docs[path]
	return ok
}

// Borrow borrows a document from the index.
// The document must not be retained outside of fn.
func (i *Index) Borrow(path string, fn func(d *Document) error) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return fn(i.docs[path].doc)
}

// Remove removes a document from the index.
func (i *Index) Remove(path string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.docs, path)
}

type document struct {
	mu  sync.Mutex
	doc *Document
}
