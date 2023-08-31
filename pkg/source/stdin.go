package source

import (
	"os"
)

const (
	stdinFilename = "(stdin)"
)

// Stdin is a source that serves standard input.
type Stdin struct {
	*VirtualFS
}

// NewStdin creates a new standard input source.
func NewStdin() (Source, error) {
	vf := NewVirtualFile(stdinFilename)

	// This reads all data until the user presses Ctrl+D.
	if _, err := vf.ReadFrom(os.Stdin); err != nil {
		return nil, err
	}

	return &Stdin{NewVirtualFS(vf)}, nil
}
