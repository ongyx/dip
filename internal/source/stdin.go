package source

import (
	"io/fs"
	"os"
)

const (
	stdinFilename = "(stdin)"
)

// Stdin is a source that reads from standard input.
type Stdin struct {
	file *VirtualFile
}

// NewStdin creates a new standard input source.
// The path is ignored.
func NewStdin(_ string) (Source, error) {
	vf := NewVirtualFile(stdinFilename)

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

	if name != stdinFilename {
		err.Err = fs.ErrNotExist
		return nil, err
	}

	return s.file, nil
}

func (s *Stdin) Root() string {
	return stdinFilename
}

func (s *Stdin) Watch(files chan<- string, errors chan<- error) {}

func (s *Stdin) Close() error {
	return nil
}
