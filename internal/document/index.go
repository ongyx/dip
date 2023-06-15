package document

import "sync"

// Index is a map of file paths to documents.
type Index struct {
	mu   sync.RWMutex
	docs map[string]*Document
}

// NewIndex creates a new index.
func NewIndex() *Index {
	return &Index{docs: make(map[string]*Document)}
}

// Add creates a new document in the index.
// If the file path exists, the existing document is returned.
func (idx *Index) Add(file string) *Document {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	d, ok := idx.docs[file]
	if !ok {
		d = &Document{}
		idx.docs[file] = d
	}

	return d
}

// Has checks if the index contains the file path.
func (idx *Index) Has(file string) bool {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	_, ok := idx.docs[file]
	return ok
}

// Get returns the document mapped to the file path.
func (idx *Index) Get(file string) *Document {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	return idx.docs[file]
}

// Remove deletes the document mapped to the file path.
func (idx *Index) Remove(file string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	delete(idx.docs, file)
}
