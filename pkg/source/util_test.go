package source

import (
	"io/fs"
	"os"
	"testing"
)

const (
	content = "# Hello World!"
)

func createfile(t *testing.T, name string, data []byte) {
	err := os.WriteFile(name, data, 0777)
	if err != nil {
		t.Errorf("failed creating file %s: %s\n", name, err)
	}
}

func tempfile(t *testing.T) *os.File {
	f, err := os.CreateTemp("", "*")
	if err != nil {
		t.Error("failed to create temp file:", err)
	}

	return f
}

func tempdir(t *testing.T) string {
	d, err := os.MkdirTemp("", "*")
	if err != nil {
		t.Error("failed to create temp dir:", err)
	}

	return d
}

func compare(t *testing.T, fsys fs.FS, name string, expected string) {
	b, err := fs.ReadFile(fsys, name)
	if err != nil {
		t.Errorf("failed to open file %s: %s\n", name, err)
	}

	got := string(b)
	if got != expected {
		t.Errorf("content in %s does not match: expected %s, got %s\n", name, expected, got)
	}
}
