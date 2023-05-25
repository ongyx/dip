package source

import (
	"io/fs"
	"os"
)

// Stdin is a source that reads from standard input.
type Stdin struct {
	file *VirtualFile
}

// NewStdin creates a new standard input source.
func NewStdin(_ string) (Source, error) {
	vf := NewVirtualFile("(stdin)")

	if _, err := vf.ReadFrom(os.Stdin); err != nil {
		return nil, err
	}

	return &Stdin{file: vf}, nil
}

func (s *Stdin) Open(name string) (fs.File, error) {
	err := &fs.PathError{Op: "open", Path: name}

	if !fs.ValidPath(name) {
		err.Err = fs.ErrInvalid
		return nil, err
	}

	if name != Root {
		err.Err = fs.ErrNotExist
		return nil, err
	}

	return s.file, nil
}

func (s *Stdin) Watch(files chan<- string, errors chan<- error) {}

func (s *Stdin) Close() error {
	return nil
}
