package document

import (
	"testing"
	"testing/fstest"
)

const (
	readme = "README.md"
)

var (
	files = fstest.MapFS{
		readme: &fstest.MapFile{
			Data: []byte("# Hello World!"),
		},
	}
)

func TestLibraryOpen(t *testing.T) {
	l := NewLibrary(files, converter)

	if _, ok := l.Open(readme); ok {
		t.Error("expected document to be non-existent")
	}

	_, err := l.Create(readme)
	if err != nil {
		t.Error("failed to create README.md")
	}

	if _, ok := l.Open(readme); !ok {
		t.Error("expected opening to succeed")
	}
}

func TestLibraryRemove(t *testing.T) {
	l := NewLibrary(files, converter)

	if l.Remove(readme) {
		t.Error("expected document to be non-existent")
	}

	_, err := l.Create(readme)
	if err != nil {
		t.Error("failed to create README.md")
	}

	if !l.Remove(readme) {
		t.Error("expected removal to succeed")
	}
}
