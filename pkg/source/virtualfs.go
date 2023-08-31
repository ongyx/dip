package source

import (
	"io/fs"
	"sync"
)

// VirtualFS is a source that serves in-memory files.
type VirtualFS struct {
	mu    sync.RWMutex
	files map[string]*VirtualFile
}

// NewVirtualFS creates a new virtual filesystem.
// If root is non-nil, it is used as the root file.
func NewVirtualFS(root *VirtualFile) *VirtualFS {
	files := make(map[string]*VirtualFile)

	if root != nil {
		// The root file may have a different name on purpose.
		files[Root] = root
	}

	return &VirtualFS{files: files}
}

// Create creates a new virtual file.
// The existing virtual file is truncated, if any.
func (vfs *VirtualFS) Create(name string) *VirtualFile {
	vfs.mu.Lock()
	defer vfs.mu.Unlock()

	vf, ok := vfs.files[name]
	if !ok {
		vf = NewVirtualFile(name)
		vfs.files[name] = vf
	} else {
		vf.Buffer.Reset()
	}

	return vf
}

// Open opens the virtual file by name for reading.
func (vfs *VirtualFS) Open(name string) (fs.File, error) {
	vfs.mu.RLock()
	defer vfs.mu.RUnlock()

	err := &fs.PathError{Op: "open", Path: name}

	if !fs.ValidPath(name) {
		err.Err = fs.ErrInvalid
		return nil, err
	}

	vf, ok := vfs.files[name]
	if !ok {
		err.Err = fs.ErrNotExist
		return nil, err
	}

	return vf, nil
}

// Watch is a no-op.
func (vfs *VirtualFS) Watch(files chan<- string, errors chan<- error) {}

// Close frees all virtual files.
func (vfs *VirtualFS) Close() error {
	vfs.files = nil
	return nil
}
