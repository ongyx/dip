package dip

import (
	"io"
	"os"
)

// Stdin is a source that reads from standard input.
type Stdin struct {
	text []byte
}

// NewStdin creates a new standard input source.
func NewStdin() (*Stdin, error) {
	text, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	return &Stdin{text: text}, nil
}

func (s *Stdin) Title(path string) string {
	return "stdin"
}

func (s *Stdin) Read(path string) ([]byte, error) {
	if path != Root {
		return nil, ErrPathNotFound
	}

	return s.text, nil
}

func (s *Stdin) Reload(queue chan<- string) {}
