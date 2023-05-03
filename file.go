package dip

import (
	"path/filepath"
)

// File is a source that reads from a single file.
type File struct {
	*Directory

	name string
}

// NewFile creates a new file source.
// path must be absolute.
func NewFile(path string) (*File, error) {
	dir, err := NewDirectory(filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	return &File{
		Directory: dir,
		name:      filepath.Base(path),
	}, nil
}

func (f *File) Title(path string) string {
	return f.name
}

func (f *File) Read(path string) ([]byte, error) {
	if path != Root {
		return nil, ErrPathNotFound
	}

	return f.Directory.Read(f.name)
}

func (f *File) Reload(queue chan<- string) {
	// we use another queue to make sure only the file we're interested in is reloaded.
	dirQueue := make(chan string)

	go func() {
		for path := range dirQueue {
			// drop all other paths
			if path == f.name {
				queue <- Root
			}
		}
	}()

	f.Directory.Reload(dirQueue)
}
