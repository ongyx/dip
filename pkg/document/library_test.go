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

	if _, err := l.Open(readme); err != ErrDocumentNotFound {
		t.Error("expected ErrDocumentNotFound for opening before creation")
	}

	_, err := l.Create(readme)
	if err != nil {
		t.Error("failed to create README.md")
	}

	if _, err := l.Open(readme); err != nil {
		t.Errorf("expected opening to succeed, got %s\n", err)
	}
}

func TestLibraryRemove(t *testing.T) {
	l := NewLibrary(files, converter)

	if err := l.Remove(readme); err != ErrDocumentNotFound {
		t.Error("expected ErrDocumentNotFound for removal before creation")
	}

	_, err := l.Create(readme)
	if err != nil {
		t.Error("failed to create README.md")
	}

	if err := l.Remove(readme); err != nil {
		t.Errorf("expected removal to succeed, got %s\n", err)
	}
}
